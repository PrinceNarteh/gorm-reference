package db

import (
	"errors"
	"fmt"
	"log"

	"gorm-reference/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ====================================================================
// Basic Transactions
// Use transactions to ensure data consistency for multiple operations.
// ====================================================================

// TransferCredits demonstrates a basic transaction
func TransferCredits(db *gorm.DB, fromUserID, toUserID uint, amount int) error {
	// Transaction wraps operations in a database transaction
	return db.Transaction(func(tx *gorm.DB) error {
		var fromUser, toUser models.User

		// Lock the rows for update to prevent race conditions
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&fromUser, fromUserID).Error; err != nil {
			return err
		}

		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&toUser, toUserID).Error; err != nil {
			return err
		}

		// Validate sufficient balance
		if fromUser.Credits < amount {
			return errors.New("insufficient credits")
		}

		// Perform the transfer
		if err := tx.Model(&fromUser).
			Update("credits", gorm.Expr("credits - ?", amount)).Error; err != nil {
			return err
		}

		if err := tx.Model(&toUser).
			Update("credits", gorm.Expr("credits + ?", amount)).Error; err != nil {
			return err
		}

		// Return nil to commit the transaction
		// Return any error to rollback
		return nil
	})
}

// ==========================================================
// Nested Transactions with Savepoints
// Use savepoints for partial rollbacks within a transaction.
// ==========================================================

// CreateOrderWithItems demonstrates nested transactions
func CreateOrderWithItems(db *gorm.DB, order *models.Order, items []models.OrderItem) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// Create the main order
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		// Process each item with savepoints
		for i, item := range items {
			item.OrderID = order.ID

			// Nested transaction creates a savepoint
			err := tx.Transaction(func(tx2 *gorm.DB) error {
				// Update inventory
				result := tx2.Model(&models.Product{}).
					Where("id = ? AND stock >= ?", item.ProductID, item.Quantity).
					Update("stock", gorm.Expr("stock - ?", item.Quantity))

				if result.RowsAffected == 0 {
					return fmt.Errorf("insufficient stock for product %d", item.ProductID)
				}

				// Create the order item
				return tx2.Create(&item).Error
			})
			if err != nil {
				// Log the error but continue with other items
				log.Printf("Failed to process item %d: %v", i, err)
				// The savepoint is rolled back, but main transaction continues
			}
		}

		// Verify at least one item was added
		var itemCount int64
		tx.Model(&models.OrderItem{}).Where("order_id = ?", order.ID).Count(&itemCount)
		if itemCount == 0 {
			return errors.New("no items could be added to the order")
		}

		return nil
	})
}
