package models

import "gorm.io/gorm"

// =======================================================================================
// Has One Relationship
// Define a has-one relationship where a model owns exactly one instance of another model.
// =======================================================================================

// Profile has a one-to-one relationship with User
type Profile struct {
	gorm.Model

	// Foreign key to User
	UserID uint `gorm:"uniqueIndex;not null" json:"userId"` // uniqueIndex ensures one profile per user

	Bio       string `gorm:"type:text" json:"bio"`
	AvatarURL string `gorm:"type:varchar(500)" json:"avatarURL"`
	Website   string `gorm:"type:varchar(255)" json:"website"`
	Location  string `gorm:"type:varchar(100)" json:"location"`

	// Social links stored as JSON
	SocialLinks map[string]string `gorm:"type:jsonb" json:"socialLinks"`
}
