package services

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"pocketploy/internal/config"
	"pocketploy/internal/models"
	"pocketploy/internal/repositories"
	"pocketploy/internal/utils"

	"github.com/google/uuid"
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo  *repositories.UserRepository
	tokenRepo *repositories.TokenRepository
	config    *config.Config
}

// NewAuthService creates a new authentication service
func NewAuthService(userRepo *repositories.UserRepository, tokenRepo *repositories.TokenRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		config:    cfg,
	}
}

// SignupParams contains parameters for user registration
type SignupParams struct {
	Username string
	Email    string
	Password string
	Request  *http.Request // HTTP request for extracting IP and User-Agent
}

// LoginParams contains parameters for user login
type LoginParams struct {
	Email    string
	Password string
	Request  *http.Request // HTTP request for extracting IP and User-Agent
}

// TokenPair contains access and refresh tokens
type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}

// RegisterUser creates a new user account
func (s *AuthService) RegisterUser(params SignupParams) (*models.User, *TokenPair, error) {
	// Normalize inputs
	params.Username = strings.ToLower(strings.TrimSpace(params.Username))
	params.Email = strings.ToLower(strings.TrimSpace(params.Email))

	// Validate username format
	if err := utils.ValidateStruct(models.SignupRequest{
		Username: params.Username,
		Email:    params.Email,
		Password: params.Password,
	}); err != nil {
		return nil, nil, fmt.Errorf("validation failed: %w", err)
	}

	// Check if username exists
	exists, err := s.userRepo.ExistsByUsername(params.Username)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to check username: %w", err)
	}
	if exists {
		return nil, nil, fmt.Errorf("username already exists")
	}

	// Check if email exists
	exists, err = s.userRepo.ExistsByEmail(params.Email)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to check email: %w", err)
	}
	if exists {
		return nil, nil, fmt.Errorf("email already exists")
	}

	// Hash password
	fmt.Printf("[DEBUG] Hashing password with bcrypt cost: %d\n", s.config.BcryptCost)
	passwordHash, err := utils.HashPassword(params.Password, s.config.BcryptCost)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to hash password: %w", err)
	}
	fmt.Printf("[DEBUG] Password hashed successfully (hash length: %d)\n", len(passwordHash))

	// Create user model
	now := time.Now().UTC()
	user := &models.User{
		ID:           uuid.New().String(),
		Username:     params.Username,
		Email:        params.Email,
		PasswordHash: passwordHash,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// Save user to database
	if err := s.userRepo.Create(user); err != nil {
		return nil, nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate tokens with request context for IP/UserAgent
	tokens, err := s.generateTokenPair(user.ID, user.Username, user.Email, params.Request)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return user, tokens, nil
}

// AuthenticateUser validates credentials and returns user with tokens
func (s *AuthService) AuthenticateUser(params LoginParams) (*models.User, *TokenPair, error) {
	// Normalize email
	params.Email = strings.ToLower(strings.TrimSpace(params.Email))

	fmt.Printf("[DEBUG] Login attempt for email: %s\n", params.Email)

	// Get user by email
	user, err := s.userRepo.GetByEmail(params.Email)
	if err != nil {
		fmt.Printf("[DEBUG] Failed to get user by email: %v\n", err)
		return nil, nil, fmt.Errorf("invalid email or password")
	}

	fmt.Printf("[DEBUG] Found user: id=%s, username=%s, is_active=%v\n", user.ID, user.Username, user.IsActive)

	// Check if user is active
	if !user.IsActive {
		fmt.Printf("[DEBUG] User account is inactive\n")
		return nil, nil, fmt.Errorf("account is inactive")
	}

	// Verify password
	fmt.Printf("[DEBUG] Verifying password (hash length: %d)\n", len(user.PasswordHash))
	if err := utils.CheckPassword(params.Password, user.PasswordHash); err != nil {
		fmt.Printf("[DEBUG] Password verification failed: %v\n", err)
		return nil, nil, fmt.Errorf("invalid email or password")
	}

	fmt.Printf("[DEBUG] Password verified successfully\n")

	// Update last login timestamp
	if err := s.userRepo.UpdateLastLogin(user.ID); err != nil {
		// Log error but don't fail the login
		fmt.Printf("Warning: failed to update last login: %v\n", err)
	}

	// Generate tokens with request context for IP/UserAgent
	tokens, err := s.generateTokenPair(user.ID, user.Username, user.Email, params.Request)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return user, tokens, nil
}

// RefreshAccessToken generates a new access token using a refresh token
func (s *AuthService) RefreshAccessToken(refreshTokenString string) (string, time.Time, error) {
	// Hash the token to look up in database
	tokenHash := utils.HashRefreshToken(refreshTokenString)

	// Get refresh token from database
	token, err := s.tokenRepo.GetByTokenHash(tokenHash)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("invalid or expired refresh token")
	}

	// Get user
	user, err := s.userRepo.GetByID(token.UserID)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("user not found")
	}

	// Check if user is active
	if !user.IsActive {
		return "", time.Time{}, fmt.Errorf("account is inactive")
	}

	// Generate new access token
	accessExpiry, _ := utils.ParseDuration(s.config.JWTAccessExpiry)
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Username, user.Email, s.config.JWTAccessSecret, accessExpiry)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to generate access token: %w", err)
	}

	expiresAt := time.Now().UTC().Add(accessExpiry)
	return accessToken, expiresAt, nil
}

// RevokeRefreshToken revokes a refresh token
func (s *AuthService) RevokeRefreshToken(refreshTokenString string) error {
	// Hash the token
	tokenHash := utils.HashRefreshToken(refreshTokenString)

	// Revoke the token
	if err := s.tokenRepo.Revoke(tokenHash); err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}

	return nil
}

// RevokeAllUserTokens revokes all refresh tokens for a user
func (s *AuthService) RevokeAllUserTokens(userID string) error {
	if err := s.tokenRepo.RevokeAllForUser(userID); err != nil {
		return fmt.Errorf("failed to revoke all tokens: %w", err)
	}
	return nil
}

// GetCurrentUser retrieves a user by ID
func (s *AuthService) GetCurrentUser(userID string) (*models.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	if !user.IsActive {
		return nil, fmt.Errorf("account is inactive")
	}

	return user, nil
}

// generateTokenPair generates both access and refresh tokens
func (s *AuthService) generateTokenPair(userID, username, email string, r *http.Request) (*TokenPair, error) {
	// Generate access token
	accessExpiry, _ := utils.ParseDuration(s.config.JWTAccessExpiry)
	accessToken, err := utils.GenerateAccessToken(userID, username, email, s.config.JWTAccessSecret, accessExpiry)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token
	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Hash token for storage
	tokenHash := utils.HashRefreshToken(refreshToken)

	// Calculate expiry
	refreshExpiry, _ := utils.ParseDuration(s.config.JWTRefreshExpiry)
	expiresAt := time.Now().UTC().Add(refreshExpiry)

	// Extract metadata from request (if available)
	var ipAddress string
	var userAgent string
	if r != nil {
		ipAddress = extractIPAddress(r)
		userAgent = r.Header.Get("User-Agent")
	}

	// Store refresh token in database
	token := &models.RefreshToken{
		ID:        uuid.New().String(),
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now().UTC(),
		IPAddress: ipAddress,
		UserAgent: userAgent,
	}

	if err := s.tokenRepo.Create(token); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// extractIPAddress extracts the client IP address from the request
func extractIPAddress(r *http.Request) string {
	// Check X-Forwarded-For header first (proxy)
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		// Take the first IP in the chain
		if idx := strings.Index(forwarded, ","); idx != -1 {
			return strings.TrimSpace(forwarded[:idx])
		}
		return strings.TrimSpace(forwarded)
	}

	// Get from RemoteAddr
	ipAddress := r.RemoteAddr

	// Strip port from IP address
	if colonIndex := strings.LastIndex(ipAddress, ":"); colonIndex != -1 {
		// Check if it's IPv6 (contains multiple colons)
		if strings.Count(ipAddress, ":") > 1 {
			// IPv6 address - remove the port after the last ]
			if bracketIndex := strings.LastIndex(ipAddress, "]"); bracketIndex != -1 {
				ipAddress = ipAddress[:bracketIndex+1]
				// Remove brackets for PostgreSQL INET type
				ipAddress = strings.Trim(ipAddress, "[]")
			}
		} else {
			// IPv4 address
			ipAddress = ipAddress[:colonIndex]
		}
	}

	return ipAddress
}
