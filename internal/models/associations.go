package models

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// =================================================
// Working with Associations
// Create, append, replace, and delete associations.
// =================================================

// AssociationOperations demonstrates working with relationships
func AssociationOperations(db *gorm.DB) {
	ctx := context.Background()

	// Create a post with tags in a single operation
	post := Post{
		Title:   "GORM Relationships",
		Content: "Understanding GORM relationships...",
		UserID:  1,
		Tags: []Tag{
			{Name: "golang", Slug: "golang"},
			{Name: "database", Slug: "database"},
		},
	}
	db.WithContext(ctx).Create(&post)

	// Append new tags to an existing post
	var existingPost Post
	db.First(&existingPost, 1)

	newTag := Tag{Name: "tutorial", Slug: "tutorial"}
	db.WithContext(ctx).Model(&existingPost).Association("Tags").Append(&newTag)

	// Replace all tags for a post
	replacementTags := []Tag{
		{Name: "go", Slug: "go"},
	}
	db.WithContext(ctx).Model(&existingPost).Association("Tags").Replace(&replacementTags)

	// Remove a specific tag from a post
	var tagToRemove Tag
	db.Where("slug = ?", "go").First(&tagToRemove)
	db.WithContext(ctx).Model(&existingPost).Association("Tags").Delete(&tagToRemove)

	// Clear all tags from a post
	db.WithContext(ctx).Model(&existingPost).Association("Tags").Clear()

	// Count associations
	count := db.WithContext(ctx).Model(&existingPost).Association("Tags").Count()
	fmt.Printf("Post has %d tags\n", count)
}
