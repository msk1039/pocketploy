package services

import (
	"fmt"
	"time"

	"pocketploy/internal/config"
	"pocketploy/internal/repositories"
)

// TokenService handles refresh token management business logic
type TokenService struct {
	tokenRepo *repositories.TokenRepository
	config    *config.Config
}

// NewTokenService creates a new token service
func NewTokenService(tokenRepo *repositories.TokenRepository, cfg *config.Config) *TokenService {
	return &TokenService{
		tokenRepo: tokenRepo,
		config:    cfg,
	}
}

// CleanupExpiredTokens removes all expired refresh tokens from the database
func (s *TokenService) CleanupExpiredTokens() (int64, error) {
	deletedCount, err := s.tokenRepo.DeleteExpired()
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup expired tokens: %w", err)
	}
	return deletedCount, nil
}

// CleanupRevokedTokens removes all revoked refresh tokens from the database
func (s *TokenService) CleanupRevokedTokens() (int64, error) {
	deletedCount, err := s.tokenRepo.DeleteRevoked()
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup revoked tokens: %w", err)
	}
	return deletedCount, nil
}

// GetUserActiveSessions returns the count of active sessions for a user
func (s *TokenService) GetUserActiveSessions(userID string) (int, error) {
	count, err := s.tokenRepo.CountByUserID(userID)
	if err != nil {
		return 0, fmt.Errorf("failed to count user sessions: %w", err)
	}
	return count, nil
}

// GetTotalActiveSessions returns the total count of active sessions across all users
func (s *TokenService) GetTotalActiveSessions() (int, error) {
	count, err := s.tokenRepo.Count()
	if err != nil {
		return 0, fmt.Errorf("failed to count total sessions: %w", err)
	}
	return count, nil
}

// RevokeAllUserSessions revokes all active sessions for a user (e.g., on password change)
func (s *TokenService) RevokeAllUserSessions(userID string) error {
	if err := s.tokenRepo.RevokeAllForUser(userID); err != nil {
		return fmt.Errorf("failed to revoke all user sessions: %w", err)
	}
	return nil
}

// RevokeSession revokes a specific session by token hash
func (s *TokenService) RevokeSession(tokenHash string) error {
	if err := s.tokenRepo.Revoke(tokenHash); err != nil {
		return fmt.Errorf("failed to revoke session: %w", err)
	}
	return nil
}

// GetUserTokens retrieves all tokens (active and inactive) for a user
func (s *TokenService) GetUserTokens(userID string) ([]TokenInfo, error) {
	tokens, err := s.tokenRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user tokens: %w", err)
	}

	tokenInfos := make([]TokenInfo, len(tokens))
	for i, token := range tokens {
		tokenInfos[i] = TokenInfo{
			ID:        token.ID,
			CreatedAt: token.CreatedAt,
			ExpiresAt: token.ExpiresAt,
			RevokedAt: token.RevokedAt,
			IPAddress: token.IPAddress,
			UserAgent: token.UserAgent,
			IsActive:  token.RevokedAt == nil && token.ExpiresAt.After(time.Now().UTC()),
			IsExpired: token.ExpiresAt.Before(time.Now().UTC()),
		}
	}

	return tokenInfos, nil
}

// TokenInfo represents display information about a token
type TokenInfo struct {
	ID        string     `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	ExpiresAt time.Time  `json:"expires_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
	IPAddress string     `json:"ip_address"`
	UserAgent string     `json:"user_agent"`
	IsActive  bool       `json:"is_active"`
	IsExpired bool       `json:"is_expired"`
}
