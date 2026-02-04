package db

import (
	"context"
	"errors"
	"time"

	"gorm-reference/internal/models"
)

// =========================================================
// Context Usage
// Always pass context for timeout and cancellation support.
// =========================================================

// FindWithTimeout demonstrates context usage for timeouts
func (r *UserRepository) FindWithTimeout(parentCtx context.Context, id uint) (*models.User, error) {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(parentCtx, 5*time.Second)
	defer cancel()

	var user models.User
	result := r.db.WithContext(ctx).First(&user, id)

	if result.Error != nil {
		if errors.Is(result.Error, context.DeadlineExceeded) {
			return nil, errors.New("query timed out")
		}
		return nil, result.Error
	}

	return &user, nil
}
