package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"pocketploy/internal/database"
)

// HealthHandler handles health check endpoints
type HealthHandler struct {
	db *database.DB
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *database.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// Health returns the API health status
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"status":    "ok",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// HealthDB checks database connection health
func (h *HealthHandler) HealthDB(w http.ResponseWriter, r *http.Request) {
	if err := h.db.Ping(); err != nil {
		respondWithJSON(w, http.StatusServiceUnavailable, map[string]interface{}{
			"status":    "error",
			"message":   "Database connection failed",
			"error":     err.Error(),
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"status":    "ok",
		"message":   "Database connection successful",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]interface{}{
		"success": false,
		"error":   message,
	})
}
