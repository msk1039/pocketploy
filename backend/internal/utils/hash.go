package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes a password using bcrypt
func HashPassword(password string, cost int) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPassword compares a password with a hash
func CheckPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
