package db

import (
	"gorm-reference/internal/models"

	"gorm.io/gorm"
)

// ==========================================================
// Batch Processing
// Process large datasets efficiently using batch operations.
// ==========================================================

// ProcessUsersInBatches processes users without loading all into memory
func ProcessUsersInBatches(db *gorm.DB, batchSize int, processor func([]models.User) error) error {
	var users []models.User

	// FindInBatches processes records in batches
	result := db.Model(&models.User{}).
		Where("is_active = ?", true).
		FindInBatches(&users, batchSize, func(tx *gorm.DB, batch int) error {
			// Process each batch
			if err := processor(users); err != nil {
				return err
			}

			// Return nil to continue to next batch
			return nil
		})

	return result.Error
}

// Rows returns an iterator for memory-efficient processing
func ProcessUsersOneByOne(db *gorm.DB, processor func(*models.User) error) error {
	rows, err := db.Model(&models.User{}).Where("is_active = ?", true).Rows()
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		// ScanRows scans a single row into the struct
		if err := db.ScanRows(rows, &user); err != nil {
			return err
		}
		if err := processor(&user); err != nil {
			return err
		}
	}

	return rows.Err()
}
