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
	StoragePath    string     `db:"storage_path" json:"storage_path"`
	CreatedAt      time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at" json:"updated_at"`
	LastAccessedAt *time.Time `db:"last_accessed_at" json:"last_accessed_at,omitempty"`
}

// InstanceStatus represents the possible states of an instance
const (
	InstanceStatusPending = "pending"
	InstanceStatusRunning = "running"
	InstanceStatusStopped = "stopped"
	InstanceStatusFailed  = "failed"
	InstanceStatusDeleted = "deleted"
)

// CreateInstanceParams holds parameters for creating a new instance
type CreateInstanceParams struct {
	UserID        uuid.UUID
	Name          string
	Slug          string
	Subdomain     string
	ContainerID   *string
	ContainerName *string
	Status        string
	StoragePath   string
}

// Create creates a new instance in the database
func (i *Instance) Create(ctx context.Context, db *sqlx.DB, params CreateInstanceParams) error {
	query := `
		INSERT INTO instances (
			user_id, name, slug, subdomain, container_id, container_name, 
			status, storage_path, created_at, updated_at
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
		params.StoragePath,
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
	i.StoragePath = params.StoragePath

	return nil
}

// FindByID retrieves an instance by its ID
func FindInstanceByID(ctx context.Context, db *sqlx.DB, id uuid.UUID) (*Instance, error) {
	var instance Instance
	query := `
		SELECT id, user_id, name, slug, subdomain, container_id, container_name,
		       status, storage_path, created_at, updated_at, last_accessed_at
		FROM instances
		WHERE id = $1 AND status != $2
	`

	err := db.GetContext(ctx, &instance, query, id, InstanceStatusDeleted)
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
		       status, storage_path, created_at, updated_at, last_accessed_at
		FROM instances
		WHERE user_id = $1 AND status != $2
		ORDER BY created_at DESC
	`

	err := db.SelectContext(ctx, &instances, query, userID, InstanceStatusDeleted)
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
		       status, storage_path, created_at, updated_at, last_accessed_at
		FROM instances
		WHERE subdomain = $1 AND status != $2
	`

	err := db.GetContext(ctx, &instance, query, subdomain, InstanceStatusDeleted)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("instance not found")
		}
		return nil, fmt.Errorf("failed to find instance: %w", err)
	}

	return &instance, nil
}

// CountUserInstances counts the number of active instances for a user
func CountUserInstances(ctx context.Context, db *sqlx.DB, userID uuid.UUID) (int, error) {
	var count int
	query := `
		SELECT COUNT(*) 
		FROM instances 
		WHERE user_id = $1 AND status NOT IN ($2, $3)
	`

	err := db.GetContext(ctx, &count, query, userID, InstanceStatusDeleted, InstanceStatusFailed)
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

// Delete marks an instance as deleted (soft delete)
func (i *Instance) Delete(ctx context.Context, db *sqlx.DB) error {
	query := `
		UPDATE instances 
		SET status = $1, updated_at = NOW()
		WHERE id = $2
	`

	result, err := db.ExecContext(ctx, query, InstanceStatusDeleted, i.ID)
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

	i.Status = InstanceStatusDeleted
	i.UpdatedAt = time.Now().UTC()

	return nil
}
