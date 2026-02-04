package db

import (
	"errors"
	"log"
	"time"

	"gorm-reference/internal/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ===========================================================
// Model Hooks
// Use hooks to add logic before or after database operations.
// ===========================================================

// BeforeCreate hook runs before inserting a new record
func (u *models.User) BeforeCreate(tx *gorm.DB) error {
	// Generate UUID if not set
	if u.UUID == "" {
		u.UUID = uuid.New().String()
	}

	// Validate email format
	if !isValidEmail(u.Email) {
		return errors.New("invalid email format")
	}

	// Hash password if provided in plain text
	if u.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.PasswordHash = string(hash)
		u.Password = "" // Clear plain text password
	}

	return nil
}

// AfterCreate hook runs after inserting a new record
func (u *User) AfterCreate(tx *gorm.DB) error {
	// Send welcome email asynchronously
	go sendWelcomeEmail(u.Email)

	// Create audit log
	return tx.Create(&AuditLog{
		Action:    "user_created",
		EntityID:  u.ID,
		Timestamp: time.Now(),
	}).Error
}

// BeforeUpdate hook runs before updating a record
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	// Track changes for audit
	if tx.Statement.Changed("Email") {
		// Log email change
		log.Printf("User %d changing email", u.ID)
	}
	return nil
}

// AfterFind hook runs after querying records
func (u *User) AfterFind(tx *gorm.DB) error {
	// Mask sensitive data
	u.PasswordHash = "[REDACTED]"
	return nil
}

// BeforeDelete hook runs before deleting a record
func (u *User) BeforeDelete(tx *gorm.DB) error {
	// Prevent deletion of admin users
	if u.Role == "admin" {
		return errors.New("cannot delete admin users")
	}
	return nil
}
