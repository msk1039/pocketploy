package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// GenerateRefreshToken generates a secure random refresh token
func GenerateRefreshToken() (string, error) {
	// Generate 32 random bytes
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// Encode as hex string (64 characters)
	return hex.EncodeToString(bytes), nil
}

// HashRefreshToken hashes a refresh token using SHA-256
func HashRefreshToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
