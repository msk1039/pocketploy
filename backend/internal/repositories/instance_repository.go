package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"pocketploy/internal/database"
	"pocketploy/internal/models"
)

// InstanceRepository handles all database operations for instances
type InstanceRepository struct {
	db *database.DB
}

// NewInstanceRepository creates a new instance repository
func NewInstanceRepository(db *database.DB) *InstanceRepository {
	return &InstanceRepository{db: db}
}

// Create inserts a new instance into the database
func (r *InstanceRepository) Create(instance *models.Instance) error {
	query := `
		INSERT INTO instances (
			id, user_id, name, slug, subdomain, container_id, container_name,
			status, data_path, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.db.Exec(query,
		instance.ID,
		instance.UserID,
		instance.Name,
		instance.Slug,
		instance.Subdomain,
		instance.ContainerID,
		instance.ContainerName,
		instance.Status,
		instance.DataPath,
		instance.CreatedAt,
		instance.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create instance: %w", err)
	}
	return nil
}

// GetByID retrieves an instance by its ID
func (r *InstanceRepository) GetByID(id string) (*models.Instance, error) {
	var instance models.Instance
	query := `SELECT * FROM instances WHERE id = $1`
	err := r.db.Get(&instance, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("instance not found")
		}
		return nil, fmt.Errorf("failed to get instance: %w", err)
	}
	return &instance, nil
}

// GetByUserID retrieves all instances for a specific user
func (r *InstanceRepository) GetByUserID(userID string) ([]*models.Instance, error) {
	var instances []*models.Instance
	query := `
		SELECT * FROM instances 
		WHERE user_id = $1 
		ORDER BY created_at DESC
	`
	err := r.db.Select(&instances, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get instances for user: %w", err)
	}
	return instances, nil
}

// GetByName retrieves an instance by name and user ID
func (r *InstanceRepository) GetByName(userID, name string) (*models.Instance, error) {
	var instance models.Instance
	query := `SELECT * FROM instances WHERE user_id = $1 AND name = $2`
	err := r.db.Get(&instance, query, userID, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("instance not found")
		}
		return nil, fmt.Errorf("failed to get instance by name: %w", err)
	}
	return &instance, nil
}

// GetByContainerID retrieves an instance by its container ID
func (r *InstanceRepository) GetByContainerID(containerID string) (*models.Instance, error) {
	var instance models.Instance
	query := `SELECT * FROM instances WHERE container_id = $1`
	err := r.db.Get(&instance, query, containerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("instance not found")
		}
		return nil, fmt.Errorf("failed to get instance by container id: %w", err)
	}
	return &instance, nil
}

// Update updates an existing instance
func (r *InstanceRepository) Update(instance *models.Instance) error {
	instance.UpdatedAt = time.Now().UTC()
	query := `
		UPDATE instances 
		SET name = $1, slug = $2, status = $3, container_id = $4, 
		    container_name = $5, subdomain = $6, data_path = $7, updated_at = $8
		WHERE id = $9
	`
	result, err := r.db.Exec(query,
		instance.Name,
		instance.Slug,
		instance.Status,
		instance.ContainerID,
		instance.ContainerName,
		instance.Subdomain,
		instance.DataPath,
		instance.UpdatedAt,
		instance.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update instance: %w", err)
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

// UpdateStatus updates only the status of an instance
func (r *InstanceRepository) UpdateStatus(id, status string) error {
	query := `UPDATE instances SET status = $1, updated_at = $2 WHERE id = $3`
	result, err := r.db.Exec(query, status, time.Now().UTC(), id)
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

	return nil
}

// UpdateContainerID updates only the container ID of an instance
func (r *InstanceRepository) UpdateContainerID(id, containerID string) error {
	query := `UPDATE instances SET container_id = $1, updated_at = $2 WHERE id = $3`
	result, err := r.db.Exec(query, containerID, time.Now().UTC(), id)
	if err != nil {
		return fmt.Errorf("failed to update container id: %w", err)
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

// Delete permanently removes an instance from the database
func (r *InstanceRepository) Delete(id string) error {
	query := `DELETE FROM instances WHERE id = $1`
	result, err := r.db.Exec(query, id)
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

// Archive moves an instance to the archive table
func (r *InstanceRepository) Archive(instance *models.Instance) error {
	// Begin transaction
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert into archive table
	archiveQuery := `
		INSERT INTO instances_archive (
			id, user_id, name, slug, subdomain, container_id, container_name,
			original_status, data_path, created_at, updated_at, deleted_at,
			deleted_by_user_id, deletion_reason, data_available, 
			data_retained_until, data_size_mb, original_subdomain
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
	`
	deletedAt := time.Now().UTC()
	dataRetainedUntil := deletedAt.AddDate(0, 0, 30) // 30 days retention

	_, err = tx.Exec(archiveQuery,
		instance.ID,
		instance.UserID,
		instance.Name,
		instance.Slug,
		instance.Subdomain,
		instance.ContainerID,
		instance.ContainerName,
		instance.Status,
		instance.DataPath,
		instance.CreatedAt,
		instance.UpdatedAt,
		deletedAt,
		instance.UserID, // deleted_by_user_id
		"Archived via repository",
		true, // data_available
		dataRetainedUntil,
		0, // data_size_mb - calculate if needed
		instance.Subdomain,
	)
	if err != nil {
		return fmt.Errorf("failed to archive instance: %w", err)
	}

	// Delete from instances table
	deleteQuery := `DELETE FROM instances WHERE id = $1`
	_, err = tx.Exec(deleteQuery, instance.ID)
	if err != nil {
		return fmt.Errorf("failed to delete instance after archive: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// ExistsByName checks if an instance with the given name exists for a user
func (r *InstanceRepository) ExistsByName(userID, name string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM instances WHERE user_id = $1 AND name = $2`
	err := r.db.QueryRow(query, userID, name).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check instance name existence: %w", err)
	}
	return count > 0, nil
}

// List retrieves all instances (admin function)
func (r *InstanceRepository) List() ([]*models.Instance, error) {
	var instances []*models.Instance
	query := `SELECT * FROM instances ORDER BY created_at DESC`
	err := r.db.Select(&instances, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list instances: %w", err)
	}
	return instances, nil
}

// Count returns the total number of instances
func (r *InstanceRepository) Count() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM instances`
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count instances: %w", err)
	}
	return count, nil
}

// CountByUserID returns the total number of instances for a user
func (r *InstanceRepository) CountByUserID(userID string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM instances WHERE user_id = $1`
	err := r.db.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count instances for user: %w", err)
	}
	return count, nil
}

// CountByStatus returns the number of instances with a specific status
func (r *InstanceRepository) CountByStatus(status string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM instances WHERE status = $1`
	err := r.db.QueryRow(query, status).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count instances by status: %w", err)
	}
	return count, nil
}

// GetByStatus retrieves all instances with a specific status
func (r *InstanceRepository) GetByStatus(status string) ([]*models.Instance, error) {
	var instances []*models.Instance
	query := `SELECT * FROM instances WHERE status = $1 ORDER BY created_at DESC`
	err := r.db.Select(&instances, query, status)
	if err != nil {
		return nil, fmt.Errorf("failed to get instances by status: %w", err)
	}
	return instances, nil
}
