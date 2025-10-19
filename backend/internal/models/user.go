package models

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID           string     `db:"id" json:"id"`
	Username     string     `db:"username" json:"username"`
	Email        string     `db:"email" json:"email"`
	PasswordHash string     `db:"password_hash" json:"-"`
	IsActive     bool       `db:"is_active" json:"is_active"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at" json:"updated_at"`
	LastLoginAt  *time.Time `db:"last_login_at" json:"last_login_at,omitempty"`
}

// SignupRequest represents the request body for user registration
type SignupRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50,alphanum_hyphen"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,password_strength"`
}

// LoginRequest represents the request body for user login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UpdateUserRequest represents the request body for updating user profile
type UpdateUserRequest struct {
	Username string `json:"username,omitempty" validate:"omitempty,min=3,max=50,alphanum_hyphen"`
	Email    string `json:"email,omitempty" validate:"omitempty,email"`
}

// UserResponse represents the public user data returned to clients
type UserResponse struct {
	ID          string     `json:"id"`
	Username    string     `json:"username"`
	Email       string     `json:"email"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
}

// ToResponse converts User to UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:          u.ID,
		Username:    u.Username,
		Email:       u.Email,
		IsActive:    u.IsActive,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
		LastLoginAt: u.LastLoginAt,
	}
}
