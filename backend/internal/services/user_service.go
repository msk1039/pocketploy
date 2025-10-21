package services

import (
	"fmt"
	"strings"

	"pocketploy/internal/config"
	"pocketploy/internal/models"
	"pocketploy/internal/repositories"
	"pocketploy/internal/utils"
)

// UserService handles user management business logic
type UserService struct {
	userRepo *repositories.UserRepository
	config   *config.Config
}

// NewUserService creates a new user service
func NewUserService(userRepo *repositories.UserRepository, cfg *config.Config) *UserService {
	return &UserService{
		userRepo: userRepo,
		config:   cfg,
	}
}

// UpdateProfileParams contains parameters for updating user profile
type UpdateProfileParams struct {
	Username *string
	Email    *string
}

// UpdatePasswordParams contains parameters for updating user password
type UpdatePasswordParams struct {
	CurrentPassword string
	NewPassword     string
}

// GetUserProfile retrieves a user's profile by ID
func (s *UserService) GetUserProfile(userID string) (*models.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	if !user.IsActive {
		return nil, fmt.Errorf("account is inactive")
	}

	return user, nil
}

// UpdateUserProfile updates a user's profile information
func (s *UserService) UpdateUserProfile(userID string, params UpdateProfileParams) (*models.User, error) {
	// Get current user
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	if !user.IsActive {
		return nil, fmt.Errorf("account is inactive")
	}

	// Update fields if provided
	updated := false

	if params.Username != nil {
		newUsername := strings.ToLower(strings.TrimSpace(*params.Username))
		if newUsername != user.Username {
			// Check if username is already taken
			exists, err := s.userRepo.ExistsByUsername(newUsername)
			if err != nil {
				return nil, fmt.Errorf("failed to check username: %w", err)
			}
			if exists {
				return nil, fmt.Errorf("username already exists")
			}
			user.Username = newUsername
			updated = true
		}
	}

	if params.Email != nil {
		newEmail := strings.ToLower(strings.TrimSpace(*params.Email))
		if newEmail != user.Email {
			// Check if email is already taken
			exists, err := s.userRepo.ExistsByEmail(newEmail)
			if err != nil {
				return nil, fmt.Errorf("failed to check email: %w", err)
			}
			if exists {
				return nil, fmt.Errorf("email already exists")
			}
			user.Email = newEmail
			updated = true
		}
	}

	// Save if anything changed
	if updated {
		if err := s.userRepo.Update(user); err != nil {
			return nil, fmt.Errorf("failed to update user: %w", err)
		}
	}

	return user, nil
}

// UpdateUserPassword updates a user's password
func (s *UserService) UpdateUserPassword(userID string, params UpdatePasswordParams) error {
	// Get current user
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	if !user.IsActive {
		return fmt.Errorf("account is inactive")
	}

	// Verify current password
	if err := utils.CheckPassword(params.CurrentPassword, user.PasswordHash); err != nil {
		return fmt.Errorf("current password is incorrect")
	}

	// Validate new password
	if len(params.NewPassword) < 8 {
		return fmt.Errorf("new password must be at least 8 characters long")
	}

	// Hash new password
	newPasswordHash, err := utils.HashPassword(params.NewPassword, s.config.BcryptCost)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	// Update password
	user.PasswordHash = newPasswordHash
	if err := s.userRepo.Update(user); err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

// DeactivateUser soft deletes a user account
func (s *UserService) DeactivateUser(userID string) error {
	// Get current user
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("user not found")
	}

	if !user.IsActive {
		return fmt.Errorf("account is already inactive")
	}

	// Soft delete (set is_active to false)
	if err := s.userRepo.Delete(userID); err != nil {
		return fmt.Errorf("failed to deactivate user: %w", err)
	}

	return nil
}

// GetUserByEmail retrieves a user by email (admin function)
func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

// GetUserByUsername retrieves a user by username (admin function)
func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	username = strings.ToLower(strings.TrimSpace(username))
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

// ListUsers retrieves all active users (admin function)
func (s *UserService) ListUsers() ([]*models.User, error) {
	users, err := s.userRepo.List()
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	return users, nil
}

// GetTotalUsers returns the total count of active users (admin function)
func (s *UserService) GetTotalUsers() (int, error) {
	count, err := s.userRepo.Count()
	if err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}
	return count, nil
}
