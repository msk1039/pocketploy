package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"pocketploy/internal/database"
	"pocketploy/internal/models"
)

// TokenRepository handles all database operations for refresh tokens
type TokenRepository struct {
	db *database.DB
}

// NewTokenRepository creates a new token repository
func NewTokenRepository(db *database.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

// Create inserts a new refresh token into the database
func (r *TokenRepository) Create(token *models.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, created_at, ip_address, user_agent)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(query,
		token.ID,
		token.UserID,
		token.TokenHash,
		token.ExpiresAt,
		token.CreatedAt,
		token.IPAddress,
		token.UserAgent,
	)
	if err != nil {
		return fmt.Errorf("failed to create refresh token: %w", err)
	}
	return nil
}

// GetByTokenHash retrieves a refresh token by its hash
func (r *TokenRepository) GetByTokenHash(tokenHash string) (*models.RefreshToken, error) {
	var token models.RefreshToken
	query := `
		SELECT * FROM refresh_tokens 
		WHERE token_hash = $1 AND revoked_at IS NULL AND expires_at > $2
	`
	err := r.db.Get(&token, query, tokenHash, time.Now().UTC())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("refresh token not found or expired")
		}
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}
	return &token, nil
}

// GetByID retrieves a refresh token by its ID
func (r *TokenRepository) GetByID(id string) (*models.RefreshToken, error) {
	var token models.RefreshToken
	query := `SELECT * FROM refresh_tokens WHERE id = $1`
	err := r.db.Get(&token, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("refresh token not found")
		}
		return nil, fmt.Errorf("failed to get refresh token by id: %w", err)
	}
	return &token, nil
}

// GetByUserID retrieves all refresh tokens for a user
func (r *TokenRepository) GetByUserID(userID string) ([]*models.RefreshToken, error) {
	var tokens []*models.RefreshToken
	query := `
		SELECT * FROM refresh_tokens 
		WHERE user_id = $1 
		ORDER BY created_at DESC
	`
	err := r.db.Select(&tokens, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tokens for user: %w", err)
	}
	return tokens, nil
}

// GetActiveByUserID retrieves all active (non-revoked, non-expired) tokens for a user
func (r *TokenRepository) GetActiveByUserID(userID string) ([]*models.RefreshToken, error) {
	var tokens []*models.RefreshToken
	query := `
		SELECT * FROM refresh_tokens 
		WHERE user_id = $1 
		AND revoked_at IS NULL 
		AND expires_at > $2
		ORDER BY created_at DESC
	`
	err := r.db.Select(&tokens, query, userID, time.Now().UTC())
	if err != nil {
		return nil, fmt.Errorf("failed to get active tokens for user: %w", err)
	}
	return tokens, nil
}

// Revoke marks a refresh token as revoked
func (r *TokenRepository) Revoke(tokenHash string) error {
	now := time.Now().UTC()
	query := `UPDATE refresh_tokens SET revoked_at = $1 WHERE token_hash = $2 AND revoked_at IS NULL`
	result, err := r.db.Exec(query, now, tokenHash)
	if err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("token not found or already revoked")
	}

	return nil
}

// RevokeByID marks a refresh token as revoked by its ID
func (r *TokenRepository) RevokeByID(id string) error {
	now := time.Now().UTC()
	query := `UPDATE refresh_tokens SET revoked_at = $1 WHERE id = $2 AND revoked_at IS NULL`
	result, err := r.db.Exec(query, now, id)
	if err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("token not found or already revoked")
	}

	return nil
}

// RevokeAllForUser revokes all tokens for a specific user
func (r *TokenRepository) RevokeAllForUser(userID string) error {
	now := time.Now().UTC()
	query := `UPDATE refresh_tokens SET revoked_at = $1 WHERE user_id = $2 AND revoked_at IS NULL`
	_, err := r.db.Exec(query, now, userID)
	if err != nil {
		return fmt.Errorf("failed to revoke all tokens for user: %w", err)
	}
	return nil
}

// DeleteExpired permanently removes expired tokens from the database
func (r *TokenRepository) DeleteExpired() (int64, error) {
	query := `DELETE FROM refresh_tokens WHERE expires_at < $1`
	result, err := r.db.Exec(query, time.Now().UTC())
	if err != nil {
		return 0, fmt.Errorf("failed to delete expired tokens: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rows, nil
}

// DeleteRevoked permanently removes revoked tokens from the database
func (r *TokenRepository) DeleteRevoked() (int64, error) {
	query := `DELETE FROM refresh_tokens WHERE revoked_at IS NOT NULL`
	result, err := r.db.Exec(query)
	if err != nil {
		return 0, fmt.Errorf("failed to delete revoked tokens: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rows, nil
}

// Count returns the total number of active tokens
func (r *TokenRepository) Count() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM refresh_tokens WHERE revoked_at IS NULL AND expires_at > $1`
	err := r.db.QueryRow(query, time.Now().UTC()).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count tokens: %w", err)
	}
	return count, nil
}

// CountByUserID returns the total number of active tokens for a user
func (r *TokenRepository) CountByUserID(userID string) (int, error) {
	var count int
	query := `
		SELECT COUNT(*) FROM refresh_tokens 
		WHERE user_id = $1 
		AND revoked_at IS NULL 
		AND expires_at > $2
	`
	err := r.db.QueryRow(query, userID, time.Now().UTC()).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count tokens for user: %w", err)
	}
	return count, nil
}
