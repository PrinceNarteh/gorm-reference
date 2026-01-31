package models

import "gorm.io/gorm"

// =============================================================================
// Belongs To Relationship
// Define a belongs-to relationship where a model references another model's ID.
// =============================================================================

// Post belongs to a User (many-to-one relationship)
type Post struct {
	gorm.Model
	Title   string `gorm:"type:varchar(255);not null" json:"title"`
	Content string `gorm:"type:text" json:"content"`

	// Foreign key field
	UserID uint `gorm:"not null;index" json:"userId"`
	// Association reference - GORM will populate this when preloading
	User User `gorm:"foreignKey:UserID" json:"user"`

	// Self-referential belongs-to for reply threads
	ParentID *uint `gorm:"index" json:"parentId"`
	Parent   *Post `gorm:"foreignKey:ParentID" json:"parent"`

	// Many-to-many with Tags
	// GORM automatically creates the join table 'post_tags'
	Tags []Tag `gorm:"many2many:post_tags;"`
}
