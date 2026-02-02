package db

import (
	"fmt"

	"gorm-reference/internal/models"

	"gorm.io/gorm"
)

// ==================================================================================
// Auto Migration
// Use AutoMigrate for development environments to automatically sync schema changes.
// ==================================================================================

// AutoMigrate creates or updates tables based on model definitions
func AutoMigrate(db *gorm.DB) error {
	// AutoMigrate will:
	// - Create tables if they don't exist
	// - Add missing columns
	// - Add missing indexes
	// It will NOT:
	// - Delete unused columns (data safety)
	// - Change column types
	// - Delete unused indexes

	err := db.AutoMigrate(
		&models.User{},
		&models.Profile{},
		&models.Post{},
		&models.Comment{},
		&models.Tag{},
	)
	if err != nil {
		return fmt.Errorf("auto migration failed: %w", err)
	}

	return nil
}
