package router

import (
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"pocketploy/internal/config"
	"pocketploy/internal/database"
	"pocketploy/internal/docker"
	appHandlers "pocketploy/internal/handlers"
	"pocketploy/internal/middleware"
	"pocketploy/internal/services"
)

// New creates a new router with all routes configured
func New(cfg *config.Config, db *database.DB, dockerClient *docker.Client) http.Handler {
	r := mux.NewRouter()

	// Initialize services
	instanceService := services.NewInstanceService(db.DB, dockerClient, cfg)

	// Initialize handlers
	healthHandler := appHandlers.NewHealthHandler(db)
	authHandler := appHandlers.NewAuthHandler(cfg, db)
	userHandler := appHandlers.NewUserHandler(db)
	instanceHandler := appHandlers.NewInstanceHandler(instanceService)

	// Health check routes (no auth required)
	r.HandleFunc("/health", healthHandler.Health).Methods("GET")
	r.HandleFunc("/health/db", healthHandler.HealthDB).Methods("GET")

	// API v1 routes
	api := r.PathPrefix("/api/v1").Subrouter()

	// Auth routes (no auth required)
	auth := api.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/signup", authHandler.Signup).Methods("POST")
	auth.HandleFunc("/login", authHandler.Login).Methods("POST")
	auth.HandleFunc("/refresh", authHandler.Refresh).Methods("POST")

	// Protected auth routes
	authProtected := api.PathPrefix("/auth").Subrouter()
	authProtected.Use(middleware.Auth(cfg))
	authProtected.HandleFunc("/logout", authHandler.Logout).Methods("POST")
	authProtected.HandleFunc("/me", authHandler.Me).Methods("GET")

	// User routes (auth required)
	users := api.PathPrefix("/users").Subrouter()
	users.Use(middleware.Auth(cfg))
	users.HandleFunc("/me", userHandler.GetMe).Methods("GET")
	users.HandleFunc("/me", userHandler.UpdateMe).Methods("PATCH")

	// Instance routes (auth required)
	instances := api.PathPrefix("/instances").Subrouter()
	instances.Use(middleware.Auth(cfg))
	instances.HandleFunc("", instanceHandler.CreateInstance).Methods("POST")
	instances.HandleFunc("", instanceHandler.ListInstances).Methods("GET")
	instances.HandleFunc("/{id}", instanceHandler.GetInstance).Methods("GET")
	instances.HandleFunc("/{id}", instanceHandler.DeleteInstance).Methods("DELETE")

	// Apply logging middleware
	loggedRouter := middleware.Logging(r)

	// Apply CORS middleware
	corsRouter := handlers.CORS(
		handlers.AllowedOrigins([]string{cfg.AllowedOrigins}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		handlers.AllowCredentials(),
		handlers.MaxAge(int((12 * time.Hour).Seconds())),
	)(loggedRouter)

	return corsRouter
}
