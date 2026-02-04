package db

import (
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// ==========================================================
// Error Handling
// Properly handle and wrap GORM errors for better debugging.
// ==========================================================

// Custom error types for domain-specific errors
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrDuplicateEmail    = errors.New("email already exists")
	ErrDuplicateUsername = errors.New("username already exists")
)

// HandleGORMError converts GORM errors to domain errors
func HandleGORMError(err error) error {
	if err == nil {
		return nil
	}

	// Check for record not found
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrUserNotFound
	}

	// Check for duplicate key violation (PostgreSQL)
	if strings.Contains(err.Error(), "duplicate key") {
		if strings.Contains(err.Error(), "email") {
			return ErrDuplicateEmail
		}
		if strings.Contains(err.Error(), "username") {
			return ErrDuplicateUsername
		}
	}

	// Return wrapped error for unexpected errors
	return fmt.Errorf("database error: %w", err)
}
