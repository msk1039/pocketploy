package services

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"pocketploy/internal/config"
	"pocketploy/internal/docker"
	"pocketploy/internal/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// InstanceService handles business logic for PocketBase instances
type InstanceService struct {
	db           *sqlx.DB
	dockerClient *docker.Client
	config       *config.Config
}

// NewInstanceService creates a new instance service
func NewInstanceService(db *sqlx.DB, dockerClient *docker.Client, cfg *config.Config) *InstanceService {
	return &InstanceService{
		db:           db,
		dockerClient: dockerClient,
		config:       cfg,
	}
}

// CreateInstanceRequest represents the request to create a new instance
type CreateInstanceRequest struct {
	UserID        uuid.UUID
	Username      string
	Name          string
	AdminEmail    string
	AdminPassword string
}

// CreateInstanceResponse represents the response after creating an instance
type CreateInstanceResponse struct {
	Instance *models.Instance
	URL      string
}

// CreateInstance creates a new PocketBase instance for a user
func (s *InstanceService) CreateInstance(ctx context.Context, req CreateInstanceRequest) (*CreateInstanceResponse, error) {
	// Validate instance name
	if err := s.validateInstanceName(req.Name); err != nil {
		return nil, err
	}

	// Check if user has reached the maximum number of instances
	count, err := models.CountUserInstances(ctx, s.db, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to count user instances: %w", err)
	}

	if count >= s.config.MaxInstancesPerUser {
		return nil, fmt.Errorf("maximum number of instances reached (%d)", s.config.MaxInstancesPerUser)
	}

	// Generate slug from instance name
	slug := s.generateSlug(req.Name)

	// Generate subdomain
	subdomain := s.generateSubdomain(req.Username, slug)

	// Check if subdomain already exists
	existing, _ := models.FindInstanceBySubdomain(ctx, s.db, subdomain)
	if existing != nil {
		return nil, fmt.Errorf("instance with this name already exists")
	}

	// Generate container name
	containerName := s.generateContainerName(req.Username, slug)

	// Generate storage path
	storagePath := s.generateStoragePath(req.Username, slug)

	// Create instance in database with creating status
	instance := &models.Instance{}
	err = instance.Create(ctx, s.db, models.CreateInstanceParams{
		UserID:        req.UserID,
		Name:          req.Name,
		Slug:          slug,
		Subdomain:     subdomain,
		ContainerID:   nil,
		ContainerName: &containerName,
		Status:        models.InstanceStatusCreating,
		DataPath:      storagePath,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create instance in database: %w", err)
	}

	// Create Docker container
	containerID, err := s.dockerClient.CreatePocketBaseContainer(ctx, docker.ContainerConfig{
		ContainerName: containerName,
		Subdomain:     subdomain,
		StoragePath:   storagePath,
		Username:      req.Username,
		InstanceSlug:  slug,
		AdminEmail:    req.AdminEmail,
		AdminPassword: req.AdminPassword,
	})

	if err != nil {
		// If container creation fails, update instance status to failed
		_ = instance.UpdateStatus(ctx, s.db, models.InstanceStatusFailed)
		return nil, fmt.Errorf("failed to create container: %w", err)
	}

	// Update instance with container ID and set status to running
	err = instance.UpdateContainerInfo(ctx, s.db, containerID, containerName)
	if err != nil {
		// Try to clean up container
		_ = s.dockerClient.RemoveContainer(ctx, containerID)
		_ = instance.UpdateStatus(ctx, s.db, models.InstanceStatusFailed)
		return nil, fmt.Errorf("failed to update instance with container info: %w", err)
	}

	// Update status to running
	err = instance.UpdateStatus(ctx, s.db, models.InstanceStatusRunning)
	if err != nil {
		return nil, fmt.Errorf("failed to update instance status: %w", err)
	}

	// Generate the full URL based on environment
	protocol := "http"
	if s.config.Env == "production" {
		protocol = "https"
	}
	url := fmt.Sprintf("%s://%s", protocol, subdomain)

	return &CreateInstanceResponse{
		Instance: instance,
		URL:      url,
	}, nil
}

// ListUserInstances retrieves all instances for a user
func (s *InstanceService) ListUserInstances(ctx context.Context, userID uuid.UUID) ([]models.Instance, error) {
	instances, err := models.FindInstancesByUserID(ctx, s.db, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list user instances: %w", err)
	}

	return instances, nil
}

// GetInstance retrieves a specific instance by ID
func (s *InstanceService) GetInstance(ctx context.Context, instanceID, userID uuid.UUID) (*models.Instance, error) {
	instance, err := models.FindInstanceByID(ctx, s.db, instanceID)
	if err != nil {
		return nil, err
	}

	// Verify the instance belongs to the user
	if instance.UserID != userID {
		return nil, fmt.Errorf("instance not found")
	}

	// Update last accessed timestamp
	_ = instance.UpdateLastAccessed(ctx, s.db)

	return instance, nil
}

// DeleteInstance archives an instance and removes its container (keeps data for 30 days)
func (s *InstanceService) DeleteInstance(ctx context.Context, instanceID, userID uuid.UUID) error {
	// Get the instance
	instance, err := models.FindInstanceByID(ctx, s.db, instanceID)
	if err != nil {
		return err
	}

	// Verify the instance belongs to the user
	if instance.UserID != userID {
		return fmt.Errorf("instance not found")
	}

	// Calculate data directory size for metadata
	dataSizeMB := 0
	if instance.DataPath != "" {
		if info, err := os.Stat(instance.DataPath); err == nil {
			dataSizeMB = int(info.Size() / 1024 / 1024) // Convert to MB
		}
	}

	// Archive the instance (moves to instances_archive table)
	_, err = models.ArchiveInstance(ctx, s.db, models.ArchiveInstanceParams{
		Instance:          instance,
		DeletedByUserID:   userID,
		DeletionReason:    "manual",
		DataSizeMB:        dataSizeMB,
		DataRetentionDays: 30, // Keep data for 30 days
	})
	if err != nil {
		return fmt.Errorf("failed to archive instance: %w", err)
	}

	// Stop and remove the container if it exists
	if instance.ContainerID != nil && *instance.ContainerID != "" {
		// Stop the container
		err = s.dockerClient.StopContainer(ctx, *instance.ContainerID)
		if err != nil {
			// Log error but continue with deletion
			fmt.Printf("Warning: failed to stop container %s: %v\n", *instance.ContainerID, err)
		}

		// Remove the container
		err = s.dockerClient.RemoveContainer(ctx, *instance.ContainerID)
		if err != nil {
			// Log error but continue with deletion
			fmt.Printf("Warning: failed to remove container %s: %v\n", *instance.ContainerID, err)
		}
	}

	// Delete instance from main table (it's now in archive)
	err = instance.Delete(ctx, s.db)
	if err != nil {
		return fmt.Errorf("failed to delete instance from main table: %w", err)
	}

	// Keep data folder for 30 days (don't delete yet)
	// A background job will clean up expired data based on data_retained_until
	fmt.Printf("Instance archived: %s (data retained until %s)\n",
		instance.Name,
		time.Now().AddDate(0, 0, 30).Format("2006-01-02"))

	return nil
}

// GetInstanceLogs retrieves logs from an instance's container
func (s *InstanceService) GetInstanceLogs(ctx context.Context, instanceID, userID uuid.UUID, tail string) (string, error) {
	instance, err := s.GetInstance(ctx, instanceID, userID)
	if err != nil {
		return "", err
	}

	if instance.ContainerID == nil || *instance.ContainerID == "" {
		return "", fmt.Errorf("instance has no container")
	}

	logs, err := s.dockerClient.GetContainerLogs(ctx, *instance.ContainerID, tail)
	if err != nil {
		return "", fmt.Errorf("failed to get container logs: %w", err)
	}

	return logs, nil
}

// GetInstanceStats retrieves statistics for an instance
func (s *InstanceService) GetInstanceStats(ctx context.Context, instanceID, userID uuid.UUID) (*docker.ContainerStats, error) {
	instance, err := s.GetInstance(ctx, instanceID, userID)
	if err != nil {
		return nil, err
	}

	if instance.ContainerID == nil || *instance.ContainerID == "" {
		return nil, fmt.Errorf("instance has no container")
	}

	stats, err := s.dockerClient.GetContainerStats(ctx, *instance.ContainerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get container stats: %w", err)
	}

	return stats, nil
}

// StartInstance starts a stopped instance
func (s *InstanceService) StartInstance(ctx context.Context, instanceID, userID uuid.UUID) error {
	instance, err := s.GetInstance(ctx, instanceID, userID)
	if err != nil {
		return err
	}

	if instance.ContainerID == nil || *instance.ContainerID == "" {
		return fmt.Errorf("instance has no container")
	}

	if instance.Status == models.InstanceStatusRunning {
		return fmt.Errorf("instance is already running")
	}

	err = s.dockerClient.StartContainer(ctx, *instance.ContainerID)
	if err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

	// Update status
	err = instance.UpdateStatus(ctx, s.db, models.InstanceStatusRunning)
	if err != nil {
		return fmt.Errorf("failed to update instance status: %w", err)
	}

	return nil
}

// StopInstance stops a running instance
func (s *InstanceService) StopInstance(ctx context.Context, instanceID, userID uuid.UUID) error {
	instance, err := s.GetInstance(ctx, instanceID, userID)
	if err != nil {
		return err
	}

	if instance.ContainerID == nil || *instance.ContainerID == "" {
		return fmt.Errorf("instance has no container")
	}

	if instance.Status == models.InstanceStatusStopped {
		return fmt.Errorf("instance is already stopped")
	}

	err = s.dockerClient.StopContainer(ctx, *instance.ContainerID)
	if err != nil {
		return fmt.Errorf("failed to stop container: %w", err)
	}

	// Update status
	err = instance.UpdateStatus(ctx, s.db, models.InstanceStatusStopped)
	if err != nil {
		return fmt.Errorf("failed to update instance status: %w", err)
	}

	return nil
}

// RestartInstance restarts an instance
func (s *InstanceService) RestartInstance(ctx context.Context, instanceID, userID uuid.UUID) error {
	instance, err := s.GetInstance(ctx, instanceID, userID)
	if err != nil {
		return err
	}

	if instance.ContainerID == nil || *instance.ContainerID == "" {
		return fmt.Errorf("instance has no container")
	}

	err = s.dockerClient.RestartContainer(ctx, *instance.ContainerID)
	if err != nil {
		return fmt.Errorf("failed to restart container: %w", err)
	}

	// Update status
	err = instance.UpdateStatus(ctx, s.db, models.InstanceStatusRunning)
	if err != nil {
		return fmt.Errorf("failed to update instance status: %w", err)
	}

	return nil
}

// validateInstanceName validates the instance name
func (s *InstanceService) validateInstanceName(name string) error {
	if len(name) < 3 || len(name) > 100 {
		return fmt.Errorf("instance name must be between 3 and 100 characters")
	}

	// Allow alphanumeric, spaces, hyphens, and underscores
	validName := regexp.MustCompile(`^[a-zA-Z0-9\s\-_]+$`)
	if !validName.MatchString(name) {
		return fmt.Errorf("instance name can only contain letters, numbers, spaces, hyphens, and underscores")
	}

	return nil
}

// generateSlug creates a URL-safe slug from the instance name
func (s *InstanceService) generateSlug(name string) string {
	// Convert to lowercase
	slug := strings.ToLower(name)

	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")

	// Replace underscores with hyphens
	slug = strings.ReplaceAll(slug, "_", "-")

	// Remove any characters that are not alphanumeric or hyphens
	reg := regexp.MustCompile(`[^a-z0-9-]+`)
	slug = reg.ReplaceAllString(slug, "")

	// Remove multiple consecutive hyphens
	reg = regexp.MustCompile(`-+`)
	slug = reg.ReplaceAllString(slug, "-")

	// Trim hyphens from start and end
	slug = strings.Trim(slug, "-")

	return slug
}

// generateSubdomain creates the full subdomain for the instance
func (s *InstanceService) generateSubdomain(username, slug string) string {
	return fmt.Sprintf("%s-%s.%s", username, slug, s.config.BaseDomain)
}

// generateContainerName creates a unique container name
func (s *InstanceService) generateContainerName(username, slug string) string {
	return fmt.Sprintf("pb-%s-%s", username, slug)
}

// generateStoragePath creates the storage path for the instance
func (s *InstanceService) generateStoragePath(username, slug string) string {
	return filepath.Join(s.config.InstancesBasePath, username, slug)
}
