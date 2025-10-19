package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"pocketploy/internal/config"
	"pocketploy/internal/utils"
)

type contextKey string

const UserIDKey contextKey = "user_id"
const UserClaimsKey contextKey = "user_claims"

// Auth middleware validates JWT token and adds user ID to context
func Auth(cfg *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				respondWithError(w, http.StatusUnauthorized, "Authorization header required")
				return
			}

			// Check if it's a Bearer token
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				respondWithError(w, http.StatusUnauthorized, "Invalid authorization header format")
				return
			}

			tokenString := parts[1]

			// Validate token
			claims, err := utils.ValidateAccessToken(tokenString, cfg.JWTAccessSecret)
			if err != nil {
				respondWithError(w, http.StatusUnauthorized, "Invalid or expired token")
				return
			}

			// Add user ID and full claims to context
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, UserClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID extracts user ID from request context
func GetUserID(r *http.Request) (string, bool) {
	userID, ok := r.Context().Value(UserIDKey).(string)
	return userID, ok
}

// GetUserClaims extracts full user claims from request context
func GetUserClaims(r *http.Request) (*utils.Claims, bool) {
	claims, ok := r.Context().Value(UserClaimsKey).(*utils.Claims)
	return claims, ok
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"error":   message,
	})
}
