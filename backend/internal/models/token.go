package models

import (
	"time"
)

// RefreshToken represents a refresh token in the system
type RefreshToken struct {
	ID        string     `db:"id" json:"id"`
	UserID    string     `db:"user_id" json:"user_id"`
	TokenHash string     `db:"token_hash" json:"-"`
	ExpiresAt time.Time  `db:"expires_at" json:"expires_at"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	RevokedAt *time.Time `db:"revoked_at" json:"revoked_at,omitempty"`
	IPAddress string     `db:"ip_address" json:"ip_address"`
	UserAgent string     `db:"user_agent" json:"user_agent"`
}

// RefreshRequest represents the request body for refreshing access token
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// LogoutRequest represents the request body for logout
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
