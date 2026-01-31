package models

import (
	"time"

	"gorm.io/gorm"
)

// =====================================================
// Many to Many Relationship
// Define a many-to-many relationship with a join table.
// =====================================================

// Tag can be associated with many Posts and vice versa
type Tag struct {
	gorm.Model
	Name string `gorm:"type:varchar(50);uniqueIndex;not null" json:"name"`
	Slug string `gorm:"type:varchar(50);uniqueIndex;not null" json:"slug"`

	// Many-to-many with Posts
	Posts []Post `gorm:"many2many:post_tags;" json:"posts"`
}

// PostTag custom join table with additional fields
type PostTag struct {
	PostID    uint      `gorm:"primaryKey" json:"postId"`
	TagID     uint      `gorm:"primaryKey" json:"tagId"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	AddedBy   uint      // Who added this tag
}

// SetupJoinTable configures the custom join table
func SetupJoinTable(db *gorm.DB) error {
	// Use SetupJoinTable to specify a custom join table model
	return db.SetupJoinTable(&Post{}, "Tags", &PostTag{})
}
