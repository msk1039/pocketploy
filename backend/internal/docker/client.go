package docker

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"pocketploy/internal/config"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

// Client wraps the Docker client with custom methods
type Client struct {
	cli    *client.Client
	config *config.Config
}

// NewClient creates a new Docker client
func NewClient(cfg *config.Config) (*Client, error) {
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithHost(cfg.DockerHost),
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client: %w", err)
	}

	// Verify Docker connection
	ctx := context.Background()
	_, err = cli.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Docker daemon: %w", err)
	}

	return &Client{
		cli:    cli,
		config: cfg,
	}, nil
}

// ContainerConfig holds configuration for creating a PocketBase container
type ContainerConfig struct {
	ContainerName string
	Subdomain     string
	StoragePath   string
	Username      string
	InstanceSlug  string
	AdminEmail    string
	AdminPassword string
}

// CreatePocketBaseContainer creates and starts a new PocketBase container with Traefik labels
func (c *Client) CreatePocketBaseContainer(ctx context.Context, cfg ContainerConfig) (string, error) {
	// Ensure storage directory exists
	if err := os.MkdirAll(cfg.StoragePath, 0755); err != nil {
		return "", fmt.Errorf("failed to create storage directory: %w", err)
	}

	// Pull the PocketBase image if not already present
	if err := c.pullImageIfNeeded(ctx); err != nil {
		return "", fmt.Errorf("failed to pull image: %w", err)
	}

	// Create entrypoint script that sets up admin and starts server
	entrypointScript := fmt.Sprintf(`#!/bin/sh
set -e
echo "Setting up PocketBase superuser..."
/usr/local/bin/pocketbase superuser upsert %s %s || true
echo "Starting PocketBase server..."
exec /usr/local/bin/pocketbase serve --http=0.0.0.0:8090
`, cfg.AdminEmail, cfg.AdminPassword)

	// Write entrypoint script to storage directory
	entrypointPath := filepath.Join(cfg.StoragePath, "entrypoint.sh")
	if err := os.WriteFile(entrypointPath, []byte(entrypointScript), 0755); err != nil {
		return "", fmt.Errorf("failed to create entrypoint script: %w", err)
	}

	// Prepare container configuration
	containerConfig := &container.Config{
		Image:      c.config.PocketBaseImage,
		Entrypoint: []string{"/pb_data/entrypoint.sh"},
		ExposedPorts: nat.PortSet{
			"8090/tcp": struct{}{},
		},
		Labels: c.buildTraefikLabels(cfg),
	}

	// Prepare host configuration with volume mount
	absStoragePath, err := filepath.Abs(cfg.StoragePath)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}

	hostConfig := &container.HostConfig{
		RestartPolicy: container.RestartPolicy{
			Name: "unless-stopped",
		},
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: absStoragePath,
				Target: "/pb_data",
			},
		},
	}

	// Network configuration
	networkConfig := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			c.config.DockerNetwork: {},
		},
	}

	// Create the container
	resp, err := c.cli.ContainerCreate(
		ctx,
		containerConfig,
		hostConfig,
		networkConfig,
		nil,
		cfg.ContainerName,
	)
	if err != nil {
		return "", fmt.Errorf("failed to create container: %w", err)
	}

	// Start the container
	if err := c.cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		// If start fails, try to remove the container
		_ = c.cli.ContainerRemove(ctx, resp.ID, container.RemoveOptions{Force: true})
		return "", fmt.Errorf("failed to start container: %w", err)
	}

	log.Printf("Created and started PocketBase container: %s (ID: %s)", cfg.ContainerName, resp.ID)
	return resp.ID, nil
}

// StopContainer stops a running container
func (c *Client) StopContainer(ctx context.Context, containerID string) error {
	timeout := 10 // seconds
	stopOptions := container.StopOptions{
		Timeout: &timeout,
	}

	if err := c.cli.ContainerStop(ctx, containerID, stopOptions); err != nil {
		return fmt.Errorf("failed to stop container: %w", err)
	}

	log.Printf("Stopped container: %s", containerID)
	return nil
}

// RemoveContainer removes a container
func (c *Client) RemoveContainer(ctx context.Context, containerID string) error {
	removeOptions := container.RemoveOptions{
		Force:         true,
		RemoveVolumes: true, // Clean up Docker volumes to save disk space
	}

	if err := c.cli.ContainerRemove(ctx, containerID, removeOptions); err != nil {
		return fmt.Errorf("failed to remove container: %w", err)
	}

	log.Printf("Removed container: %s", containerID)
	return nil
}

// ListUserContainers lists all containers for a specific user
func (c *Client) ListUserContainers(ctx context.Context, username string) ([]string, error) {
	containers, err := c.cli.ContainerList(ctx, container.ListOptions{
		All: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	var userContainers []string
	prefix := fmt.Sprintf("pb-%s-", username)

	for _, container := range containers {
		for _, name := range container.Names {
			// Docker names start with '/'
			if len(name) > 1 && name[0] == '/' {
				name = name[1:]
			}
			if len(name) >= len(prefix) && name[:len(prefix)] == prefix {
				userContainers = append(userContainers, container.ID)
			}
		}
	}

	return userContainers, nil
}

// GetContainerStatus checks if a container is running
func (c *Client) GetContainerStatus(ctx context.Context, containerID string) (string, error) {
	containerJSON, err := c.cli.ContainerInspect(ctx, containerID)
	if err != nil {
		return "", fmt.Errorf("failed to inspect container: %w", err)
	}

	if containerJSON.State.Running {
		return "running", nil
	}
	return "stopped", nil
}

// GetContainerLogs retrieves logs from a container
func (c *Client) GetContainerLogs(ctx context.Context, containerID string, tail string) (string, error) {
	options := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Tail:       tail, // e.g., "100" for last 100 lines, "all" for all logs
		Timestamps: true,
	}

	reader, err := c.cli.ContainerLogs(ctx, containerID, options)
	if err != nil {
		return "", fmt.Errorf("failed to get container logs: %w", err)
	}
	defer reader.Close()

	logs, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("failed to read logs: %w", err)
	}

	return string(logs), nil
}

// StartContainer starts a stopped container
func (c *Client) StartContainer(ctx context.Context, containerID string) error {
	if err := c.cli.ContainerStart(ctx, containerID, container.StartOptions{}); err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}

	log.Printf("Started container: %s", containerID)
	return nil
}

// RestartContainer restarts a container
func (c *Client) RestartContainer(ctx context.Context, containerID string) error {
	timeout := 10 // seconds
	stopOptions := container.StopOptions{
		Timeout: &timeout,
	}

	if err := c.cli.ContainerRestart(ctx, containerID, stopOptions); err != nil {
		return fmt.Errorf("failed to restart container: %w", err)
	}

	log.Printf("Restarted container: %s", containerID)
	return nil
}

// GetContainerStats retrieves container statistics
func (c *Client) GetContainerStats(ctx context.Context, containerID string) (*ContainerStats, error) {
	containerJSON, err := c.cli.ContainerInspect(ctx, containerID)
	if err != nil {
		return nil, fmt.Errorf("failed to inspect container: %w", err)
	}

	stats := &ContainerStats{
		ContainerID: containerID,
		Status:      "stopped",
		Health:      "unknown",
		StartedAt:   "",
		CreatedAt:   containerJSON.Created,
	}

	if containerJSON.State.Running {
		stats.Status = "running"
		stats.StartedAt = containerJSON.State.StartedAt
		stats.Health = "healthy" // PocketBase doesn't have built-in health checks
	}

	return stats, nil
}

// ContainerStats holds container statistics
type ContainerStats struct {
	ContainerID string `json:"container_id"`
	Status      string `json:"status"`
	Health      string `json:"health"`
	StartedAt   string `json:"started_at"`
	CreatedAt   string `json:"created_at"`
}

// buildTraefikLabels creates the necessary Traefik labels for routing
// In development (HTTP): simple routing without HTTPS
// In production (HTTPS): routing with TLS and redirects
func (c *Client) buildTraefikLabels(cfg ContainerConfig) map[string]string {
	routerName := cfg.ContainerName
	labels := map[string]string{
		"traefik.enable": "true",
		fmt.Sprintf("traefik.http.routers.%s.rule", routerName):                      fmt.Sprintf("Host(`%s`)", cfg.Subdomain),
		fmt.Sprintf("traefik.http.routers.%s.entrypoints", routerName):               "web",
		fmt.Sprintf("traefik.http.services.%s.loadbalancer.server.port", routerName): "8090",
		"traefik.docker.network": c.config.TraefikNetwork,
	}

	// Add HTTPS configuration only in production
	if c.config.Env == "production" {
		// Add redirect middleware for HTTP -> HTTPS
		labels[fmt.Sprintf("traefik.http.routers.%s.middlewares", routerName)] = fmt.Sprintf("%s-redirect", routerName)
		labels[fmt.Sprintf("traefik.http.middlewares.%s-redirect.redirectscheme.scheme", routerName)] = "https"

		// Add HTTPS router
		labels[fmt.Sprintf("traefik.http.routers.%s-secure.rule", routerName)] = fmt.Sprintf("Host(`%s`)", cfg.Subdomain)
		labels[fmt.Sprintf("traefik.http.routers.%s-secure.entrypoints", routerName)] = "websecure"
		labels[fmt.Sprintf("traefik.http.routers.%s-secure.tls", routerName)] = "true"
	}

	return labels
}

// pullImageIfNeeded pulls the PocketBase image if it's not already present
func (c *Client) pullImageIfNeeded(ctx context.Context) error {
	// Check if image exists
	_, _, err := c.cli.ImageInspectWithRaw(ctx, c.config.PocketBaseImage)
	if err == nil {
		// Image already exists
		return nil
	}

	// Pull the image
	log.Printf("Pulling PocketBase image: %s", c.config.PocketBaseImage)
	reader, err := c.cli.ImagePull(ctx, c.config.PocketBaseImage, image.PullOptions{})
	if err != nil {
		return fmt.Errorf("failed to pull image: %w", err)
	}
	defer reader.Close()

	// Wait for pull to complete
	_, err = io.Copy(io.Discard, reader)
	if err != nil {
		return fmt.Errorf("failed to wait for image pull: %w", err)
	}

	log.Printf("Successfully pulled image: %s", c.config.PocketBaseImage)
	return nil
}

// Close closes the Docker client connection
func (c *Client) Close() error {
	return c.cli.Close()
}
