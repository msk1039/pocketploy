package main

import (
	"fmt"
	"log"
	"os"

	"pocketploy/internal/config"
	"pocketploy/internal/database"
	"pocketploy/internal/utils"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run cmd/reset-password/main.go <email> <new_password>")
		fmt.Println("Example: go run cmd/reset-password/main.go user@example.com newpassword123")
		os.Exit(1)
	}

	email := os.Args[1]
	newPassword := os.Args[2]

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Build database DSN
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode)

	// Connect to database
	db, err := database.New(dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Hash the new password
	passwordHash, err := utils.HashPassword(newPassword, cfg.BcryptCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	// Update the password in database
	query := `UPDATE users SET password_hash = $1 WHERE email = $2`
	result, err := db.Exec(query, passwordHash, email)
	if err != nil {
		log.Fatalf("Failed to update password: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatalf("Failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		log.Fatalf("No user found with email: %s", email)
	}

	fmt.Printf("✅ Password updated successfully for user: %s\n", email)
	fmt.Printf("   Hash length: %d characters\n", len(passwordHash))
	fmt.Printf("   Bcrypt cost: %d\n", cfg.BcryptCost)
}
