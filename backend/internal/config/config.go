package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	// Server Configuration
	Port string
	Host string
	Env  string

	// Database Configuration
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// JWT Configuration
	JWTAccessSecret  string
	JWTRefreshSecret string
	JWTAccessExpiry  string
	JWTRefreshExpiry string

	// CORS Configuration
	AllowedOrigins string

	// Bcrypt Configuration
	BcryptCost int

	// Docker Configuration
	DockerHost      string
	DockerNetwork   string
	PocketBaseImage string
	TraefikNetwork  string

	// Instance Configuration
	LocalIP             string
	InstancesBasePath   string
	MaxInstancesPerUser int
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists (ignore error in production)
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	config := &Config{
		// Server Configuration
		Port: getEnv("PORT", "8080"),
		Host: getEnv("HOST", "localhost"),
		Env:  getEnv("ENV", "development"),

		// Database Configuration
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "pocketploy"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		// JWT Configuration
		JWTAccessSecret:  getEnv("JWT_ACCESS_SECRET", ""),
		JWTRefreshSecret: getEnv("JWT_REFRESH_SECRET", ""),
		JWTAccessExpiry:  getEnv("JWT_ACCESS_EXPIRY", "15m"),
		JWTRefreshExpiry: getEnv("JWT_REFRESH_EXPIRY", "168h"),

		// CORS Configuration
		AllowedOrigins: getEnv("ALLOWED_ORIGINS", "http://localhost:3000"),

		// Bcrypt Configuration
		BcryptCost: getEnvAsInt("BCRYPT_COST", 12),

		// Docker Configuration
		DockerHost:      getEnv("DOCKER_HOST", "unix:///var/run/docker.sock"),
		DockerNetwork:   getEnv("DOCKER_NETWORK", "pocketploy-network"),
		PocketBaseImage: getEnv("POCKETBASE_IMAGE", "ghcr.io/muchobien/pocketbase:latest"),
		TraefikNetwork:  getEnv("TRAEFIK_NETWORK", "pocketploy-network"),

		// Instance Configuration
		LocalIP:             getEnv("LOCAL_IP", "127.0.0.1"),
		InstancesBasePath:   getEnv("INSTANCES_BASE_PATH", "./instances"),
		MaxInstancesPerUser: getEnvAsInt("MAX_INSTANCES_PER_USER", 5),
	}

	// Validate required fields
	if err := config.validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// validate checks if all required configuration values are set
func (c *Config) validate() error {
	if c.DBPassword == "" {
		return fmt.Errorf("DB_PASSWORD is required")
	}

	if c.JWTAccessSecret == "" {
		return fmt.Errorf("JWT_ACCESS_SECRET is required")
	}

	if len(c.JWTAccessSecret) < 32 {
		return fmt.Errorf("JWT_ACCESS_SECRET must be at least 32 characters long")
	}

	if c.JWTRefreshSecret == "" {
		return fmt.Errorf("JWT_REFRESH_SECRET is required")
	}

	if len(c.JWTRefreshSecret) < 32 {
		return fmt.Errorf("JWT_REFRESH_SECRET must be at least 32 characters long")
	}

	if c.BcryptCost < 10 || c.BcryptCost > 14 {
		return fmt.Errorf("BCRYPT_COST must be between 10 and 14")
	}

	return nil
}

// GetDSN returns the PostgreSQL connection string
func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost,
		c.DBPort,
		c.DBUser,
		c.DBPassword,
		c.DBName,
		c.DBSSLMode,
	)
}

// getEnv reads an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt reads an environment variable as integer or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Warning: Invalid integer value for %s, using default: %d", key, defaultValue)
		return defaultValue
	}

	return value
}
