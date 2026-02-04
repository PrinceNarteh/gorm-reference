package db

import (
	"database/sql"
	"fmt"

	"gorm-reference/internal/models"

	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/tern/v2/migrate"

	"gorm.io/driver/postgres"
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

// RunMigrations executes database migrations from the specified path
func RunMigrations(db *sql.DB, migrationsPath string) error {
	// Create postgres driver for migrate
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", &err)
	}

	// Create migrate instance
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed: %w", err)
	}

	return nil
}

// RollbackMigration rolls back the last migration
func RollbackMigration(db *sql.DB, migrationsPath string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"postgres",
		driver,
	)
	if err != nil {
		return err
	}

	// Roll back one step
	return m.Steps(-1)
}
