package database

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// DB holds the database connection
type DB struct {
	*sqlx.DB
}

// New creates a new database connection
func New(dsn string) (*DB, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	log.Println("Successfully connected to database")

	return &DB{db}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.DB.Close()
}

// Ping checks if the database connection is alive
func (db *DB) Ping() error {
	return db.DB.Ping()
}
