package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims represents JWT claims
type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Type     string `json:"type"`
	jwt.RegisteredClaims
}

// GenerateAccessToken generates a new JWT access token
func GenerateAccessToken(userID, username, email, secret string, expiry time.Duration) (string, error) {
	now := time.Now()
	claims := &Claims{
		UserID:   userID,
		Username: username,
		Email:    email,
		Type:     "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(expiry)),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateAccessToken validates and parses a JWT access token
func ValidateAccessToken(tokenString, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		// Verify token type
		if claims.Type != "access" {
			return nil, errors.New("invalid token type")
		}
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// ParseDuration parses a duration string (e.g., "15m", "24h")
func ParseDuration(s string) (time.Duration, error) {
	return time.ParseDuration(s)
}
