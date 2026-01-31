package repository

import (
	"context"
	"errors"
	"time"

	"gorm-reference/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	_           UserRepository = (*userRepository)(nil)
	ErrNotFound                = errors.New("user not found")
)

type UserRepository interface {
	Create(context.Context, *models.User) error
	CreateBatch(context.Context, *[]models.User) error
	Upsert(context.Context, *models.User) error
	FindByID(ctx context.Context, id uint) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindAll(ctx context.Context, page, perPage int) ([]models.User, int64, error)
	FindWithFilters(ctx context.Context, filters models.UserFilters) ([]models.User, error)
	Update(ctx context.Context, id uint, updates models.User) error
	Save(ctx context.Context, user *models.User) error
	UpdateLastLogin(ctx context.Context, id uint) error
	IncreaseLoginCount(ctx context.Context, id uint) error
}

// UserRepository handles user database operations
type userRepository struct {
	db *gorm.DB
}

func (u *userRepository) userQuery(opts ...clause.Expression) gorm.Interface[models.User] {
	return gorm.G[models.User](u.db, opts...)
}

// ===================================================================================
// Create Operations
// Insert records into the database with various methods for single and batch inserts.
// ===================================================================================

// Create inserts a single user into the database
func (u *userRepository) Create(ctx context.Context, user *models.User) error {
	return u.userQuery().Create(ctx, user)
}

// CreateBatch inserts multiple users in a single query for better performance
func (u *userRepository) CreateBatch(ctx context.Context, users *[]models.User) error {
	// CreateInBatches processes records in batches to avoid memory issues
	// Second parameter is the batch size
	return u.userQuery().CreateInBatches(ctx, users, 100)
}

// Upsert creates or updates a user based on conflict columns
func (u *userRepository) Upsert(ctx context.Context, user *models.User) error {
	// Clauses for handling conflicts (upsert)
	return u.userQuery(clause.OnConflict{
		Columns:   []clause.Column{{Name: "email"}},
		DoUpdates: clause.AssignmentColumns([]string{"username", "updated_at"}),
	}).Create(ctx, user)
}

// ================================================================
// Read Operations
// Query records with various conditions, ordering, and pagination.
// ================================================================

// FindByID retrieves a user by their ID
func (u *userRepository) FindByID(ctx context.Context, id uint) (*models.User, error) {
	user, err := u.userQuery().Where("id = ?", id).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

// FindByEmail retrieves a user by their email
func (u *userRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := u.userQuery().Where("email = ?", email).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (u *userRepository) FindAll(ctx context.Context, page, perPage int) ([]models.User, int64, error) {
	// Count total records for pagination
	total, err := u.userQuery().Count(ctx, "*")
	if err != nil {
		return nil, 0, err
	}

	// calculate offset for pagination
	offset := (page - 1) * perPage

	// Retrieve paginated records
	users, err := u.userQuery().
		Order("created_at DESC").
		Offset(offset).
		Limit(perPage).
		Find(ctx)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// FindWithFilters retrieves users matching multiple conditions
func (u *userRepository) FindWithFilters(ctx context.Context, filters models.UserFilters) ([]models.User, error) {
	// Start building the query
	query := u.userQuery()

	// Apply filters conditionally
	if filters.IsActive != nil {
		query.Where("is_active = ?", *filters.IsActive)
	}
	if filters.Username != "" {
		query.Where("username ILIKE ?", "%"+filters.Username+"%")
	}
	if !filters.CreatedAfter.IsZero() {
		query.Where("created_at > ?", filters.CreatedAfter)
	}

	return query.Find(ctx)
}

// ====================================================================
// Update Operations
// Update records with various strategies for partial and full updates.
// ====================================================================

// Update updates specific fields of a user
func (u *userRepository) Update(ctx context.Context, id uint, updates models.User) error {
	// Updates only the specified fields
	rowsAffected, err := u.userQuery().Where("id = ?", id).Updates(ctx, updates)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// Save updates all fields of a user (including zero values)
func (u *userRepository) Save(ctx context.Context, user *models.User) error {
	// Save will update all fields, including zero values
	// Use this when you want to explicitly set fields to zero/empty
	result := u.db.WithContext(ctx).Save(user)
	return result.Error
}

// UpdateLastLogin updates a single column without running hooks
func (u *userRepository) UpdateLastLogin(ctx context.Context, id uint) error {
	now := time.Now()
	rowsAffected, err := u.userQuery().Where("id = ?", id).Update(ctx, "last_login_at", now)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// IncrementCounter demonstrates atomic counter updates
func (u *userRepository) IncreaseLoginCount(ctx context.Context, id uint) error {
	// Use gorm.Expr for SQL expressions
	result := u.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", id).
		Update("login_count", gorm.Expr("login_count + ?", 1))
	return result.Error
}

// =======================================================================
// Delete Operations
// Delete records with soft delete support and permanent deletion options.
// =======================================================================

// Delete performs a soft delete (sets deleted_at)
func (u *userRepository) Delete(ctx context.Context, id uint) error {
	// With gorm.Model, Delete sets deleted_at instead of removing the row
	rowsAffected, err := u.userQuery().Where("id = ?", id).Delete(ctx)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// HardDelete permanently removes a record from the database
func (u *userRepository) HardDelete(ctx context.Context, id uint) error {
	// Unscoped bypasses soft delete and permanently removes the record
	result := u.db.WithContext(ctx).Unscoped().Delete(&models.User{}, id)
	return result.Error
}

// DeleteByCondition deletes multiple records matching a condition
func (u *userRepository) DeleteInactiveUsers(ctx context.Context, before time.Time) (int, error) {
	rowsAffected, err := u.userQuery().
		Where("is_active = ? AND last_login_at < ?", false, before).
		Delete(ctx)
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}

// Restore recovers a soft-deleted record
func (u *userRepository) Restore(ctx context.Context, id uint) error {
	result := u.db.WithContext(ctx).
		Unscoped().
		Model(&models.User{}).
		Where("id = ?", id).
		Update("deleted_at", nil)
	return result.Error
}
