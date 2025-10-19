package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"pocketploy/internal/middleware"
	"pocketploy/internal/services"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// InstanceHandler handles PocketBase instance endpoints
type InstanceHandler struct {
	instanceService *services.InstanceService
}

// NewInstanceHandler creates a new instance handler
func NewInstanceHandler(instanceService *services.InstanceService) *InstanceHandler {
	return &InstanceHandler{
		instanceService: instanceService,
	}
}

// CreateInstanceRequest represents the request to create a new instance
type CreateInstanceRequest struct {
	Name string `json:"name" validate:"required,min=3,max=100"`
}

// CreateInstance handles POST /api/v1/instances
func (h *InstanceHandler) CreateInstance(w http.ResponseWriter, r *http.Request) {
	// Get user claims from context (set by auth middleware)
	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Parse user ID
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid user ID")
		return
	}

	// Parse request body
	var req CreateInstanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if req.Name == "" {
		respondWithError(w, http.StatusBadRequest, "Instance name is required")
		return
	}

	if len(req.Name) < 3 || len(req.Name) > 100 {
		respondWithError(w, http.StatusBadRequest, "Instance name must be between 3 and 100 characters")
		return
	}

	// Create instance
	result, err := h.instanceService.CreateInstance(r.Context(), services.CreateInstanceRequest{
		UserID:   userID,
		Username: claims.Username,
		Name:     req.Name,
	})

	if err != nil {
		// Log the actual error for debugging
		fmt.Printf("Error creating instance: %v\n", err)

		// Check for specific errors
		if err.Error() == "maximum number of instances reached (5)" {
			respondWithError(w, http.StatusForbidden, err.Error())
			return
		}
		if err.Error() == "instance with this name already exists" {
			respondWithError(w, http.StatusConflict, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to create instance")
		return
	}

	// Return success response
	respondWithJSON(w, http.StatusCreated, map[string]interface{}{
		"success":  true,
		"message":  "Instance created successfully",
		"instance": result.Instance,
		"url":      result.URL,
	})
}

// ListInstances handles GET /api/v1/instances
func (h *InstanceHandler) ListInstances(w http.ResponseWriter, r *http.Request) {
	// Get user claims from context
	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Parse user ID
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid user ID")
		return
	}

	// Get user's instances
	instances, err := h.instanceService.ListUserInstances(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to list instances")
		return
	}

	// Return instances
	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"success":   true,
		"instances": instances,
	})
}

// GetInstance handles GET /api/v1/instances/:id
func (h *InstanceHandler) GetInstance(w http.ResponseWriter, r *http.Request) {
	// Get user claims from context
	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Parse user ID
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid user ID")
		return
	}

	// Get instance ID from URL
	vars := mux.Vars(r)
	instanceID, err := uuid.Parse(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid instance ID")
		return
	}

	// Get instance
	instance, err := h.instanceService.GetInstance(r.Context(), instanceID, userID)
	if err != nil {
		if err.Error() == "instance not found" {
			respondWithError(w, http.StatusNotFound, "Instance not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to get instance")
		return
	}

	// Return instance
	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"success":  true,
		"instance": instance,
	})
}

// DeleteInstance handles DELETE /api/v1/instances/:id
func (h *InstanceHandler) DeleteInstance(w http.ResponseWriter, r *http.Request) {
	// Get user claims from context
	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Parse user ID
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid user ID")
		return
	}

	// Get instance ID from URL
	vars := mux.Vars(r)
	instanceID, err := uuid.Parse(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid instance ID")
		return
	}

	// Delete instance
	err = h.instanceService.DeleteInstance(r.Context(), instanceID, userID)
	if err != nil {
		if err.Error() == "instance not found" {
			respondWithError(w, http.StatusNotFound, "Instance not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to delete instance")
		return
	}

	// Return success response
	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Instance deleted successfully",
	})
}

// GetInstanceLogs retrieves logs for a specific instance
func (h *InstanceHandler) GetInstanceLogs(w http.ResponseWriter, r *http.Request) {
	// Get user claims from context
	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Parse user ID
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid user ID")
		return
	}

	// Get instance ID from URL
	vars := mux.Vars(r)
	instanceID, err := uuid.Parse(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid instance ID")
		return
	}

	// Get tail parameter (default to 100 lines)
	tail := r.URL.Query().Get("tail")
	if tail == "" {
		tail = "100"
	}

	// Get logs
	logs, err := h.instanceService.GetInstanceLogs(r.Context(), instanceID, userID, tail)
	if err != nil {
		if err.Error() == "instance not found" {
			respondWithError(w, http.StatusNotFound, "Instance not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve logs")
		return
	}

	// Return logs
	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"logs":    logs,
	})
}

// GetInstanceStats retrieves statistics for a specific instance
func (h *InstanceHandler) GetInstanceStats(w http.ResponseWriter, r *http.Request) {
	// Get user claims from context
	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Parse user ID
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid user ID")
		return
	}

	// Get instance ID from URL
	vars := mux.Vars(r)
	instanceID, err := uuid.Parse(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid instance ID")
		return
	}

	// Get stats
	stats, err := h.instanceService.GetInstanceStats(r.Context(), instanceID, userID)
	if err != nil {
		if err.Error() == "instance not found" {
			respondWithError(w, http.StatusNotFound, "Instance not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve stats")
		return
	}

	// Return stats
	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"stats":   stats,
	})
}

// StartInstance starts a stopped instance
func (h *InstanceHandler) StartInstance(w http.ResponseWriter, r *http.Request) {
	// Get user claims from context
	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Parse user ID
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid user ID")
		return
	}

	// Get instance ID from URL
	vars := mux.Vars(r)
	instanceID, err := uuid.Parse(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid instance ID")
		return
	}

	// Start instance
	err = h.instanceService.StartInstance(r.Context(), instanceID, userID)
	if err != nil {
		if err.Error() == "instance not found" {
			respondWithError(w, http.StatusNotFound, "Instance not found")
			return
		}
		if err.Error() == "instance is already running" {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to start instance")
		return
	}

	// Return success response
	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Instance started successfully",
	})
}

// StopInstance stops a running instance
func (h *InstanceHandler) StopInstance(w http.ResponseWriter, r *http.Request) {
	// Get user claims from context
	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Parse user ID
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid user ID")
		return
	}

	// Get instance ID from URL
	vars := mux.Vars(r)
	instanceID, err := uuid.Parse(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid instance ID")
		return
	}

	// Stop instance
	err = h.instanceService.StopInstance(r.Context(), instanceID, userID)
	if err != nil {
		if err.Error() == "instance not found" {
			respondWithError(w, http.StatusNotFound, "Instance not found")
			return
		}
		if err.Error() == "instance is already stopped" {
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to stop instance")
		return
	}

	// Return success response
	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Instance stopped successfully",
	})
}

// RestartInstance restarts an instance
func (h *InstanceHandler) RestartInstance(w http.ResponseWriter, r *http.Request) {
	// Get user claims from context
	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Parse user ID
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid user ID")
		return
	}

	// Get instance ID from URL
	vars := mux.Vars(r)
	instanceID, err := uuid.Parse(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid instance ID")
		return
	}

	// Restart instance
	err = h.instanceService.RestartInstance(r.Context(), instanceID, userID)
	if err != nil {
		if err.Error() == "instance not found" {
			respondWithError(w, http.StatusNotFound, "Instance not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to restart instance")
		return
	}

	// Return success response
	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Instance restarted successfully",
	})
}
