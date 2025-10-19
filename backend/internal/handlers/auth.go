package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"pocketploy/internal/config"
	"pocketploy/internal/database"
	"pocketploy/internal/middleware"
	"pocketploy/internal/models"
	"pocketploy/internal/utils"

	"github.com/google/uuid"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	cfg *config.Config
	db  *database.DB
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(cfg *config.Config, db *database.DB) *AuthHandler {
	return &AuthHandler{
		cfg: cfg,
		db:  db,
	}
}

// Signup handles user registration
func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
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

	// Normalize username and email
	req.Username = strings.ToLower(strings.TrimSpace(req.Username))
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	// Check if username already exists
	var count int
	err := h.db.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1", req.Username).Scan(&count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Database error checking username")
		return
	}
	if count > 0 {
		respondWithError(w, http.StatusConflict, "Username already exists")
		return
	}

	// Check if email already exists
	err = h.db.QueryRow("SELECT COUNT(*) FROM users WHERE email = $1", req.Email).Scan(&count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Database error checking email")
		return
	}
	if count > 0 {
		respondWithError(w, http.StatusConflict, "Email already exists")
		return
	}

	// Hash password
	passwordHash, err := utils.HashPassword(req.Password, h.cfg.BcryptCost)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	// Create user
	userID := uuid.New().String()
	now := time.Now().UTC()

	query := `
		INSERT INTO users (id, username, email, password_hash, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err = h.db.Exec(query, userID, req.Username, req.Email, passwordHash, true, now, now)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// Generate tokens
	accessExpiry, _ := utils.ParseDuration(h.cfg.JWTAccessExpiry)
	accessToken, err := utils.GenerateAccessToken(userID, req.Username, req.Email, h.cfg.JWTAccessSecret, accessExpiry)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate access token")
		return
	}

	refreshToken, _, expiresAt, err := h.createRefreshToken(userID, r)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate refresh token")
		return
	}

	user := models.User{
		ID:        userID,
		Username:  req.Username,
		Email:     req.Email,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Return response
	respondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "User created successfully",
		"data": map[string]interface{}{
			"user":          user.ToResponse(),
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"expires_at":    expiresAt,
		},
	})
}

// Login handles user authentication
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
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

	// Normalize email
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	// Get user by email
	var user models.User
	err := h.db.Get(&user, "SELECT * FROM users WHERE email = $1", req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Database error")
		return
	}

	// Check if user is active
	if !user.IsActive {
		respondWithError(w, http.StatusUnauthorized, "Account is inactive")
		return
	}

	// Verify password
	if err := utils.CheckPassword(req.Password, user.PasswordHash); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Update last login
	now := time.Now().UTC()
	_, err = h.db.Exec("UPDATE users SET last_login_at = $1, updated_at = $2 WHERE id = $3", now, now, user.ID)
	if err != nil {
		// Log error but don't fail the login
	}

	// Generate tokens
	accessExpiry, _ := utils.ParseDuration(h.cfg.JWTAccessExpiry)
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Username, user.Email, h.cfg.JWTAccessSecret, accessExpiry)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate access token")
		return
	}

	refreshToken, _, expiresAt, err := h.createRefreshToken(user.ID, r)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate refresh token")
		return
	}

	user.LastLoginAt = &now

	// Return response
	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Login successful",
		"data": map[string]interface{}{
			"user":          user.ToResponse(),
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"expires_at":    expiresAt,
		},
	})
}

// Refresh handles token refresh
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req models.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.RefreshToken == "" {
		respondWithError(w, http.StatusBadRequest, "Refresh token is required")
		return
	}

	// Hash the token to look up in database
	tokenHash := utils.HashRefreshToken(req.RefreshToken)

	// Get refresh token from database
	var token models.RefreshToken
	query := `
		SELECT * FROM refresh_tokens 
		WHERE token_hash = $1 AND revoked_at IS NULL AND expires_at > $2
	`
	err := h.db.Get(&token, query, tokenHash, time.Now().UTC())
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusUnauthorized, "Invalid or expired refresh token")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Database error")
		return
	}

	// Get user
	var user models.User
	err = h.db.Get(&user, "SELECT * FROM users WHERE id = $1", token.UserID)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "User not found")
		return
	}

	// Generate new access token
	accessExpiry, _ := utils.ParseDuration(h.cfg.JWTAccessExpiry)
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Username, user.Email, h.cfg.JWTAccessSecret, accessExpiry)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate access token")
		return
	}

	// Return response
	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"access_token": accessToken,
			"expires_at":   time.Now().UTC().Add(accessExpiry),
		},
	})
}

// Logout handles user logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req models.LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.RefreshToken == "" {
		respondWithError(w, http.StatusBadRequest, "Refresh token is required")
		return
	}

	// Hash the token
	tokenHash := utils.HashRefreshToken(req.RefreshToken)

	// Revoke the token
	now := time.Now().UTC()
	query := `UPDATE refresh_tokens SET revoked_at = $1 WHERE token_hash = $2 AND revoked_at IS NULL`
	_, err := h.db.Exec(query, now, tokenHash)
	if err != nil {
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

	// Get user from database
	var user models.User
	err := h.db.Get(&user, "SELECT * FROM users WHERE id = $1", userID)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "User not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Database error")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"user": user.ToResponse(),
		},
	})
}

// createRefreshToken creates a new refresh token and stores it in the database
func (h *AuthHandler) createRefreshToken(userID string, r *http.Request) (string, string, time.Time, error) {
	// Generate refresh token
	refreshToken, err := utils.GenerateRefreshToken()
	if err != nil {
		return "", "", time.Time{}, err
	}

	// Hash token for storage
	tokenHash := utils.HashRefreshToken(refreshToken)

	// Calculate expiry
	refreshExpiry, _ := utils.ParseDuration(h.cfg.JWTRefreshExpiry)
	expiresAt := time.Now().UTC().Add(refreshExpiry)

	// Get client IP and user agent
	ipAddress := r.RemoteAddr
	// Strip port from IP address if present (e.g., "127.0.0.1:54321" -> "127.0.0.1")
	if colonIndex := strings.LastIndex(ipAddress, ":"); colonIndex != -1 {
		ipAddress = ipAddress[:colonIndex]
	}
	userAgent := r.Header.Get("User-Agent")

	// Store in database
	tokenID := uuid.New().String()
	query := `
		INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, created_at, ip_address, user_agent)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err = h.db.Exec(query, tokenID, userID, tokenHash, expiresAt, time.Now().UTC(), ipAddress, userAgent)
	if err != nil {
		return "", "", time.Time{}, err
	}

	return refreshToken, tokenHash, expiresAt, nil
}
