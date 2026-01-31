// Package models
package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents the users table with common fields
type User struct {
	// Embed gorm.Model for ID, CreatedAt, UpdatedAt, DeletedAt
	gorm.Model

	// Use index tag for better query performance
	FirstName *string `gorm:"type:varchar(100);index" json:"firstName"`
	LastName  *string `gorm:"type:varchar(100);index" json:"lastName"`

	// Use column tag to customize the database column name
	Email *string `gorm:"column:email;type:varchar(255);uniqueIndex;not null" json:"email"`

	// Add size constraint directly in the type
	Username *string `gorm:"type:varchar(100);uniqueIndex;not null" json:"username"`

	// Password hash should never be selected by default
	PasswordHash string `gorm:"type:varchar(255);not null;->:false;<-" json:"password"`

	// Boolean field with default value
	IsActive bool `gorm:"default:true" json:"isActive"`

	// Timestamp fields with auto-update
	LastLoginAt *time.Time `gorm:"index" json:"lastLoginAt"`

	// JSON field for flexible data storage
	Preferences map[string]any `gorm:"type:jsonb" json:"preferences"`

	// Relationships (defined later)
	Posts    []Post    `gorm:"foreignKey:UserID" json:"posts"`
	Profile  *Profile  `gorm:"foreignKey:UserID" json:"profile"`
	Comments []Comment `gorm:"foreignKey:UserID" json:"comments"`
}

// TableName overrides the default table name
func (User) TableName() string {
	return "users"
}

// UserFilters contains optional filters for querying users
type UserFilters struct {
	IsActive     *bool
	Username     string
	CreatedAfter time.Time
}
