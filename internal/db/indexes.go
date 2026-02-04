package db

import (
	"gorm-reference/internal/models"

	"gorm.io/gorm"
)

// CreateIndexes manually creates indexes using GORM's Migrator
func CreateIndexes(db *gorm.DB) error {
	migrator := db.Migrator()

	// Create a composite index
	if !migrator.HasIndex(&models.User{}, "idx_users_name") {
		err := migrator.CreateIndex(&models.User{}, "idx_users_name")
		if err != nil {
			return err
		}
	}

	// Create a partial index using raw SQL
	err := db.Exec(`
        CREATE INDEX IF NOT EXISTS idx_active_users
        ON users(email)
        WHERE is_active = true AND deleted_at IS NULL
    `).Error
	if err != nil {
		return err
	}

	return nil
}

// AddConstraints adds foreign key and check constraints
func AddConstraints(db *gorm.DB) error {
	// Add foreign key constraint with cascade
	err := db.Exec(`
        ALTER TABLE posts
        ADD CONSTRAINT fk_posts_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
    `).Error
	if err != nil {
		return err
	}

	// Add check constraint
	err = db.Exec(`
        ALTER TABLE users
        ADD CONSTRAINT chk_email_format
        CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$')
    `).Error

	return err
}
