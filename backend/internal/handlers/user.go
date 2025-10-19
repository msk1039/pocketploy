package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"pocketploy/internal/database"
	"pocketploy/internal/middleware"
	"pocketploy/internal/models"
	"pocketploy/internal/utils"
)

// UserHandler handles user-related endpoints
type UserHandler struct {
	db *database.DB
}

// NewUserHandler creates a new user handler
func NewUserHandler(db *database.DB) *UserHandler {
	return &UserHandler{db: db}
}

// GetMe returns the current user's profile
func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
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

// UpdateMe updates the current user's profile
func (h *UserHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := middleware.GetUserID(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req models.UpdateUserRequest
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

	// Build update query dynamically
	updates := []string{}
	args := []interface{}{}
	argCount := 1

	if req.Username != "" {
		req.Username = strings.ToLower(strings.TrimSpace(req.Username))

		// Check if username is already taken
		var existingID string
		err := h.db.Get(&existingID, "SELECT id FROM users WHERE username = $1 AND id != $2", req.Username, userID)
		if err != sql.ErrNoRows {
			if err == nil {
				respondWithError(w, http.StatusConflict, "Username already exists")
				return
			}
		}

		updates = append(updates, "username = $"+string(rune('0'+argCount)))
		args = append(args, req.Username)
		argCount++
	}

	if req.Email != "" {
		req.Email = strings.ToLower(strings.TrimSpace(req.Email))

		// Check if email is already taken
		var existingID string
		err := h.db.Get(&existingID, "SELECT id FROM users WHERE email = $1 AND id != $2", req.Email, userID)
		if err != sql.ErrNoRows {
			if err == nil {
				respondWithError(w, http.StatusConflict, "Email already exists")
				return
			}
		}

		updates = append(updates, "email = $"+string(rune('0'+argCount)))
		args = append(args, req.Email)
		argCount++
	}

	if len(updates) == 0 {
		respondWithError(w, http.StatusBadRequest, "No fields to update")
		return
	}

	// Add updated_at
	updates = append(updates, "updated_at = $"+string(rune('0'+argCount)))
	args = append(args, time.Now().UTC())
	argCount++

	// Add user ID for WHERE clause
	args = append(args, userID)

	// Execute update
	query := "UPDATE users SET " + strings.Join(updates, ", ") + " WHERE id = $" + string(rune('0'+argCount))
	_, err := h.db.Exec(query, args...)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update user")
		return
	}

	// Get updated user
	var user models.User
	err = h.db.Get(&user, "SELECT * FROM users WHERE id = $1", userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch updated user")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Profile updated successfully",
		"data": map[string]interface{}{
			"user": user.ToResponse(),
		},
	})
}
