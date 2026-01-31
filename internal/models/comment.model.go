package models

import "gorm.io/gorm"

// Comment belongs to both User and Post
type Comment struct {
	gorm.Model
	Content string `gorm:"type:text;not null" json:"content"`

	// Foreign keys
	UserID uint `gorm:"not null;index" json:"userId"`
	PostID uint `gorm:"not null;index" json:"postId"`

	// Associations
	User User `gorm:"foreignKey:UserID" json:"user"`
	Post Post `gorm:"foreignKey:PostID" json:"post"`
}
