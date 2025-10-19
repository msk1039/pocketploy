package models

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// Instance represents a PocketBase instance
type Instance struct {
	ID             uuid.UUID  `db:"id" json:"id"`
	UserID         uuid.UUID  `db:"user_id" json:"user_id"`
	Name           string     `db:"name" json:"name"`
	Slug           string     `db:"slug" json:"slug"`
	Subdomain      string     `db:"subdomain" json:"subdomain"`
	ContainerID    *string    `db:"container_id" json:"container_id,omitempty"`
	ContainerName  *string    `db:"container_name" json:"container_name,omitempty"`
	Status         string     `db:"status" json:"status"`
	DataPath       string     `db:"data_path" json:"data_path"`
	CreatedAt      time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at" json:"updated_at"`
	LastAccessedAt *time.Time `db:"last_accessed_at" json:"last_accessed_at,omitempty"`
}

// InstanceStatus represents the possible states of an instance
const (
	InstanceStatusCreating = "creating"
	InstanceStatusRunning  = "running"
	InstanceStatusStopped  = "stopped"
	InstanceStatusFailed   = "failed"
)

// ArchivedInstance represents a deleted instance with metadata for restore capability
type ArchivedInstance struct {
	ID                uuid.UUID  `db:"id" json:"id"`
	UserID            uuid.UUID  `db:"user_id" json:"user_id"`
	Name              string     `db:"name" json:"name"`
	Slug              string     `db:"slug" json:"slug"`
	Subdomain         string     `db:"subdomain" json:"subdomain"`
	ContainerID       *string    `db:"container_id" json:"container_id,omitempty"`
	ContainerName     *string    `db:"container_name" json:"container_name,omitempty"`
	OriginalStatus    string     `db:"original_status" json:"original_status"`
	DataPath          string     `db:"data_path" json:"data_path"`
	CreatedAt         time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time  `db:"updated_at" json:"updated_at"`
	LastAccessedAt    *time.Time `db:"last_accessed_at" json:"last_accessed_at,omitempty"`
	DeletedAt         time.Time  `db:"deleted_at" json:"deleted_at"`
	DeletedByUserID   uuid.UUID  `db:"deleted_by_user_id" json:"deleted_by_user_id"`
	DeletionReason    string     `db:"deletion_reason" json:"deletion_reason"`
	DataAvailable     bool       `db:"data_available" json:"data_available"`
	DataRetainedUntil time.Time  `db:"data_retained_until" json:"data_retained_until"`
	DataSizeMB        int        `db:"data_size_mb" json:"data_size_mb"`
	OriginalSubdomain string     `db:"original_subdomain" json:"original_subdomain"`
}

// ArchiveInstanceParams holds parameters for archiving an instance
type ArchiveInstanceParams struct {
	Instance          *Instance
	DeletedByUserID   uuid.UUID
	DeletionReason    string
	DataSizeMB        int
	DataRetentionDays int // Number of days to retain data (default 30)
}

// CreateInstanceParams holds parameters for creating a new instance
type CreateInstanceParams struct {
	UserID        uuid.UUID
	Name          string
	Slug          string
	Subdomain     string
	ContainerID   *string
	ContainerName *string
	Status        string
	DataPath      string
}

// Create creates a new instance in the database
func (i *Instance) Create(ctx context.Context, db *sqlx.DB, params CreateInstanceParams) error {
	query := `
		INSERT INTO instances (
			user_id, name, slug, subdomain, container_id, container_name, 
			status, data_path, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW()
		) RETURNING id, created_at, updated_at
	`

	err := db.QueryRowxContext(
		ctx,
		query,
		params.UserID,
		params.Name,
		params.Slug,
		params.Subdomain,
		params.ContainerID,
		params.ContainerName,
		params.Status,
		params.DataPath,
	).Scan(&i.ID, &i.CreatedAt, &i.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create instance: %w", err)
	}

	// Populate the instance object
	i.UserID = params.UserID
	i.Name = params.Name
	i.Slug = params.Slug
	i.Subdomain = params.Subdomain
	i.ContainerID = params.ContainerID
	i.ContainerName = params.ContainerName
	i.Status = params.Status
	i.DataPath = params.DataPath

	return nil
}

// FindByID retrieves an instance by its ID
func FindInstanceByID(ctx context.Context, db *sqlx.DB, id uuid.UUID) (*Instance, error) {
	var instance Instance
	query := `
		SELECT id, user_id, name, slug, subdomain, container_id, container_name,
		       status, data_path, created_at, updated_at, last_accessed_at
		FROM instances
		WHERE id = $1
	`

	err := db.GetContext(ctx, &instance, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("instance not found")
		}
		return nil, fmt.Errorf("failed to find instance: %w", err)
	}

	return &instance, nil
}

// FindByUserID retrieves all instances for a user
func FindInstancesByUserID(ctx context.Context, db *sqlx.DB, userID uuid.UUID) ([]Instance, error) {
	var instances []Instance
	query := `
		SELECT id, user_id, name, slug, subdomain, container_id, container_name,
		       status, data_path, created_at, updated_at, last_accessed_at
		FROM instances
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	err := db.SelectContext(ctx, &instances, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find instances: %w", err)
	}

	return instances, nil
}

// FindBySubdomain retrieves an instance by its subdomain
func FindInstanceBySubdomain(ctx context.Context, db *sqlx.DB, subdomain string) (*Instance, error) {
	var instance Instance
	query := `
		SELECT id, user_id, name, slug, subdomain, container_id, container_name,
		       status, data_path, created_at, updated_at, last_accessed_at
		FROM instances
		WHERE subdomain = $1
	`

	err := db.GetContext(ctx, &instance, query, subdomain)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("instance not found")
		}
		return nil, fmt.Errorf("failed to find instance: %w", err)
	}

	return &instance, nil
}

// CountUserInstances counts the number of active instances for a user (excluding failed)
func CountUserInstances(ctx context.Context, db *sqlx.DB, userID uuid.UUID) (int, error) {
	var count int
	query := `
		SELECT COUNT(*) 
		FROM instances 
		WHERE user_id = $1 AND status != $2
	`

	err := db.GetContext(ctx, &count, query, userID, InstanceStatusFailed)
	if err != nil {
		return 0, fmt.Errorf("failed to count instances: %w", err)
	}

	return count, nil
}

// UpdateStatus updates the status of an instance
func (i *Instance) UpdateStatus(ctx context.Context, db *sqlx.DB, status string) error {
	query := `
		UPDATE instances 
		SET status = $1, updated_at = NOW()
		WHERE id = $2
	`

	result, err := db.ExecContext(ctx, query, status, i.ID)
	if err != nil {
		return fmt.Errorf("failed to update instance status: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("instance not found")
	}

	i.Status = status
	i.UpdatedAt = time.Now().UTC()

	return nil
}

// UpdateContainerInfo updates the container ID and name
func (i *Instance) UpdateContainerInfo(ctx context.Context, db *sqlx.DB, containerID, containerName string) error {
	query := `
		UPDATE instances 
		SET container_id = $1, container_name = $2, updated_at = NOW()
		WHERE id = $3
	`

	result, err := db.ExecContext(ctx, query, containerID, containerName, i.ID)
	if err != nil {
		return fmt.Errorf("failed to update container info: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("instance not found")
	}

	i.ContainerID = &containerID
	i.ContainerName = &containerName
	i.UpdatedAt = time.Now().UTC()

	return nil
}

// UpdateLastAccessed updates the last accessed timestamp
func (i *Instance) UpdateLastAccessed(ctx context.Context, db *sqlx.DB) error {
	query := `
		UPDATE instances 
		SET last_accessed_at = NOW()
		WHERE id = $1
	`

	_, err := db.ExecContext(ctx, query, i.ID)
	if err != nil {
		return fmt.Errorf("failed to update last accessed: %w", err)
	}

	now := time.Now().UTC()
	i.LastAccessedAt = &now

	return nil
}

// Delete permanently deletes an instance from the database (should be archived first)
func (i *Instance) Delete(ctx context.Context, db *sqlx.DB) error {
	query := `DELETE FROM instances WHERE id = $1`

	result, err := db.ExecContext(ctx, query, i.ID)
	if err != nil {
		return fmt.Errorf("failed to delete instance: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("instance not found")
	}

	return nil
}

// ArchiveInstance moves an instance to the archive table with metadata
func ArchiveInstance(ctx context.Context, db *sqlx.DB, params ArchiveInstanceParams) (*ArchivedInstance, error) {
	instance := params.Instance

	// Calculate data retention date (default 30 days)
	retentionDays := params.DataRetentionDays
	if retentionDays == 0 {
		retentionDays = 30
	}
	dataRetainedUntil := time.Now().UTC().AddDate(0, 0, retentionDays)

	archived := &ArchivedInstance{
		ID:                instance.ID,
		UserID:            instance.UserID,
		Name:              instance.Name,
		Slug:              instance.Slug,
		Subdomain:         instance.Subdomain,
		ContainerID:       instance.ContainerID,
		ContainerName:     instance.ContainerName,
		OriginalStatus:    instance.Status,
		DataPath:          instance.DataPath,
		CreatedAt:         instance.CreatedAt,
		UpdatedAt:         instance.UpdatedAt,
		LastAccessedAt:    instance.LastAccessedAt,
		DeletedAt:         time.Now().UTC(),
		DeletedByUserID:   params.DeletedByUserID,
		DeletionReason:    params.DeletionReason,
		DataAvailable:     true,
		DataRetainedUntil: dataRetainedUntil,
		DataSizeMB:        params.DataSizeMB,
		OriginalSubdomain: instance.Subdomain,
	}

	query := `
		INSERT INTO instances_archive (
			id, user_id, name, slug, subdomain, container_id, container_name,
			original_status, data_path, created_at, updated_at, last_accessed_at,
			deleted_at, deleted_by_user_id, deletion_reason, data_available,
			data_retained_until, data_size_mb, original_subdomain
		) VALUES (
			:id, :user_id, :name, :slug, :subdomain, :container_id, :container_name,
			:original_status, :data_path, :created_at, :updated_at, :last_accessed_at,
			:deleted_at, :deleted_by_user_id, :deletion_reason, :data_available,
			:data_retained_until, :data_size_mb, :original_subdomain
		)
	`

	_, err := db.NamedExecContext(ctx, query, archived)
	if err != nil {
		return nil, fmt.Errorf("failed to archive instance: %w", err)
	}

	return archived, nil
}

// FindArchivedInstancesByUserID retrieves all archived instances for a user
func FindArchivedInstancesByUserID(ctx context.Context, db *sqlx.DB, userID uuid.UUID) ([]ArchivedInstance, error) {
	var instances []ArchivedInstance
	query := `
		SELECT id, user_id, name, slug, subdomain, container_id, container_name,
		       original_status, data_path, created_at, updated_at, last_accessed_at,
		       deleted_at, deleted_by_user_id, deletion_reason, data_available,
		       data_retained_until, data_size_mb, original_subdomain
		FROM instances_archive
		WHERE user_id = $1
		ORDER BY deleted_at DESC
	`

	err := db.SelectContext(ctx, &instances, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find archived instances: %w", err)
	}

	return instances, nil
}

// FindArchivedInstanceByID retrieves a specific archived instance
func FindArchivedInstanceByID(ctx context.Context, db *sqlx.DB, id uuid.UUID, userID uuid.UUID) (*ArchivedInstance, error) {
	var archived ArchivedInstance
	query := `
		SELECT id, user_id, name, slug, subdomain, container_id, container_name,
		       original_status, data_path, created_at, updated_at, last_accessed_at,
		       deleted_at, deleted_by_user_id, deletion_reason, data_available,
		       data_retained_until, data_size_mb, original_subdomain
		FROM instances_archive
		WHERE id = $1 AND user_id = $2
	`

	err := db.GetContext(ctx, &archived, query, id, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("archived instance not found")
		}
		return nil, fmt.Errorf("failed to find archived instance: %w", err)
	}

	return &archived, nil
}

// UpdateDataAvailability updates the data_available flag for an archived instance
func UpdateArchivedDataAvailability(ctx context.Context, db *sqlx.DB, id uuid.UUID, available bool) error {
	query := `
		UPDATE instances_archive 
		SET data_available = $1
		WHERE id = $2
	`

	_, err := db.ExecContext(ctx, query, available, id)
	if err != nil {
		return fmt.Errorf("failed to update data availability: %w", err)
	}

	return nil
}

// FindExpiredArchivedInstances finds archived instances whose data retention period has expired
func FindExpiredArchivedInstances(ctx context.Context, db *sqlx.DB) ([]ArchivedInstance, error) {
	var instances []ArchivedInstance
	query := `
		SELECT id, user_id, name, slug, subdomain, container_id, container_name,
		       original_status, data_path, created_at, updated_at, last_accessed_at,
		       deleted_at, deleted_by_user_id, deletion_reason, data_available,
		       data_retained_until, data_size_mb, original_subdomain
		FROM instances_archive
		WHERE data_retained_until < NOW() AND data_available = true
		ORDER BY data_retained_until ASC
	`

	err := db.SelectContext(ctx, &instances, query)
	if err != nil {
		return nil, fmt.Errorf("failed to find expired archived instances: %w", err)
	}

	return instances, nil
}
