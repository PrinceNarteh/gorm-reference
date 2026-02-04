package repository

import (
	"context"

	"gorm-reference/internal/models"

	"gorm.io/gorm"
)

var _ PostRepository = (*postRepository)(nil)

type PostRepository interface {
	FindPostsWithDetails(ctx context.Context, page, pageSize int) ([]models.Post, error)
	FindPostsByUserEmail(ctx context.Context, email string) ([]models.Post, error)
	FindPostsWithActiveComments(ctx context.Context) ([]models.Post, error)
	FindPostsWithUserData(ctx context.Context) ([]models.Post, error)
	FindPopularPosts(ctx context.Context, minComments int) ([]models.Post, error)
}

type postRepository struct {
	db *gorm.DB
}

// ===================================================================================
// Query Optimization
// Eager Loading with Preload
// Use Preload to avoid N+1 query problems when loading associations.
// ===================================================================================

// FindPostsWithDetails demonstrates proper eager loading
func (r *postRepository) FindPostsWithDetails(ctx context.Context, page, pageSize int) ([]models.Post, error) {
	var posts []models.Post
	offset := (page - 1) * pageSize

	// Preload loads associations in separate queries
	// This avoids the N+1 problem
	result := r.db.WithContext(ctx).
		Preload("User").          // Load the post author
		Preload("Tags").          // Load all tags
		Preload("Comments").      // Load all comments
		Preload("Comments.User"). // Load comment authors (nested preload)
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&posts)

	return posts, result.Error
}

// Conditional Preload loads associations only when certain conditions are met
func (r *postRepository) FindPostsWithActiveComments(ctx context.Context) ([]models.Post, error) {
	var posts []models.Post

	result := r.db.WithContext(ctx).
		Preload("Comments", func(db *gorm.DB) *gorm.DB {
			// Only load non-deleted comments, ordered by date
			return db.Where("deleted_at IS NULL").Order("created_at DESC")
		}).
		Find(&posts)

	return posts, result.Error
}

// =============================================================
// Using Joins for Better Performance
// Use Joins when you need to filter by associated table fields.
// =============================================================

// FindPostsByUserEmail demonstrates using joins for filtering
func (r *postRepository) FindPostsByUserEmail(ctx context.Context, email string) ([]models.Post, error) {
	var posts []models.Post

	// Joins is more efficient when filtering by associated table columns
	result := r.db.WithContext(ctx).
		Joins("JOIN users ON users.id = posts.user_id").
		Where("users.email = ?", email).
		Find(&posts)

	return posts, result.Error
}

// FindPostsWithUserData loads posts with user data in a single query
func (r *postRepository) FindPostsWithUserData(ctx context.Context) ([]models.Post, error) {
	var posts []models.Post

	// Using Joins with struct population
	result := r.db.WithContext(ctx).
		Joins("User"). // Smart join that populates the User field
		Find(&posts)

	return posts, result.Error
}

// ComplexJoinQuery demonstrates multi-table joins
func (r *postRepository) FindPopularPosts(ctx context.Context, minComments int) ([]models.Post, error) {
	var posts []models.Post

	result := r.db.WithContext(ctx).
		Select("posts.*, COUNT(comments.id) as comment_count").
		Joins("LEFT JOIN comments ON comments.post_id = posts.id").
		Group("posts.id").
		Having("COUNT(comments.id) >= ?", minComments).
		Order("comment_count DESC").
		Find(&posts)

	return posts, result.Error
}

// ===========================================================================
// Select Specific Fields
// Only select fields you need to reduce memory usage and improve performance.
// ===========================================================================

// FindPostSummaries loads only necessary fields
func (r *postRepository) FindPostSummaries(ctx context.Context) ([]models.PostSummary, error) {
	var summaries []models.PostSummary

	result := r.db.WithContext(ctx).
		Model(&models.Post{}).
		Select("posts.id, posts.title, posts.created_at, users.username as user_name").
		Joins("JOIN users ON users.id = posts.user_id").
		Scan(&summaries)

	return summaries, result.Error
}

// Pluck extracts a single column into a slice
func (r *userRepository) GetAllEmails(ctx context.Context) ([]string, error) {
	var emails []string
	result := r.db.WithContext(ctx).Model(&models.User{}).Pluck("email", &emails)
	return emails, result.Error
}
