package repository

import (
	"context"
	"errors"

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
func (u *userRepository) FinndWithFilters(ctx context.Context, filters models.UserFilters) ([]models.User, error) {
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
