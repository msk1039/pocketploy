package handlers

import (
	"encoding/json"
	"net/http"

	"pocketploy/internal/middleware"
	"pocketploy/internal/models"
	"pocketploy/internal/services"
	"pocketploy/internal/utils"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authService *services.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Signup handles user registration
func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	// Parse request
	var req models.SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		validationErrors := utils.GetValidationErrors(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Validation failed",
			"details": validationErrors,
		})
		return
	}

	// Call service to create user
	user, tokens, err := h.authService.RegisterUser(services.SignupParams{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Request:  r,
	})
	if err != nil {
		// Map service errors to HTTP status codes
		statusCode := http.StatusInternalServerError
		if err.Error() == "username already exists" || err.Error() == "email already exists" {
			statusCode = http.StatusConflict
		} else if err.Error() == "validation failed" {
			statusCode = http.StatusBadRequest
		}
		respondWithError(w, statusCode, err.Error())
		return
	}

	// Return response
	respondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "User created successfully",
		"data": map[string]interface{}{
			"user":          user.ToResponse(),
			"access_token":  tokens.AccessToken,
			"refresh_token": tokens.RefreshToken,
			"expires_at":    tokens.ExpiresAt,
		},
	})
}

// Login handles user authentication
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Parse request
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := utils.ValidateStruct(req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	// Call service to authenticate user
	user, tokens, err := h.authService.AuthenticateUser(services.LoginParams{
		Email:    req.Email,
		Password: req.Password,
		Request:  r,
	})
	if err != nil {
		// Map service errors to HTTP status codes
		statusCode := http.StatusInternalServerError
		if err.Error() == "invalid email or password" || err.Error() == "account is inactive" {
			statusCode = http.StatusUnauthorized
		}
		respondWithError(w, statusCode, err.Error())
		return
	}

	// Return response
	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Login successful",
		"data": map[string]interface{}{
			"user":          user.ToResponse(),
			"access_token":  tokens.AccessToken,
			"refresh_token": tokens.RefreshToken,
			"expires_at":    tokens.ExpiresAt,
		},
	})
}

// Refresh handles token refresh
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	// Parse request
	var req models.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.RefreshToken == "" {
		respondWithError(w, http.StatusBadRequest, "Refresh token is required")
		return
	}

	// Call service to refresh access token
	accessToken, expiresAt, err := h.authService.RefreshAccessToken(req.RefreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Return response
	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"access_token": accessToken,
			"expires_at":   expiresAt,
		},
	})
}

// Logout handles user logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Parse request
	var req models.LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.RefreshToken == "" {
		respondWithError(w, http.StatusBadRequest, "Refresh token is required")
		return
	}

	// Call service to revoke token
	if err := h.authService.RevokeRefreshToken(req.RefreshToken); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to revoke token")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Logged out successfully",
	})
}

// Me returns the current user's information
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := middleware.GetUserID(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Call service to get user
	user, err := h.authService.GetCurrentUser(userID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "user not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "account is inactive" {
			statusCode = http.StatusUnauthorized
		}
		respondWithError(w, statusCode, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"user": user.ToResponse(),
		},
	})
}
