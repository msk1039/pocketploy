package handlers

import (
	"encoding/json"
	"net/http"

	"pocketploy/internal/middleware"
	"pocketploy/internal/models"
	"pocketploy/internal/services"
	"pocketploy/internal/utils"
)

// UserHandler handles user-related endpoints
type UserHandler struct {
	userService *services.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetMe returns the current user's profile
func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := middleware.GetUserID(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Call service to get user profile
	user, err := h.userService.GetUserProfile(userID)
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

// UpdateMe updates the current user's profile
func (h *UserHandler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := middleware.GetUserID(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Parse request
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

	// Check if there are any fields to update
	if req.Username == "" && req.Email == "" {
		respondWithError(w, http.StatusBadRequest, "No fields to update")
		return
	}

	// Prepare update parameters
	params := services.UpdateProfileParams{}
	if req.Username != "" {
		params.Username = &req.Username
	}
	if req.Email != "" {
		params.Email = &req.Email
	}

	// Call service to update user profile
	user, err := h.userService.UpdateUserProfile(userID, params)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "username already exists" || err.Error() == "email already exists" {
			statusCode = http.StatusConflict
		} else if err.Error() == "user not found" {
			statusCode = http.StatusNotFound
		} else if err.Error() == "account is inactive" {
			statusCode = http.StatusUnauthorized
		}
		respondWithError(w, statusCode, err.Error())
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
