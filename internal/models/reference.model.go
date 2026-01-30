package models

import "time"

// TagReference demonstrates all common GORM struct tags
type TagReference struct {
	// Primary key - auto-increment by default for uint types
	ID uint `gorm:"primaryKey"`

	// Column name customization
	Name string `gorm:"column:full_name"`

	// Type specification for database column
	Description string `gorm:"type:text"`

	// Size constraint (shorthand for varchar)
	Code string `gorm:"size:50"`

	// Unique constraint
	Email string `gorm:"unique"`

	// Not null constraint
	Username string `gorm:"not null"`

	// Default value
	Status string `gorm:"default:'pending'"`

	// Index creation
	Category string `gorm:"index"`

	// Named index
	Tag string `gorm:"index:idx_tag"`

	// Composite index (multiple fields with same index name)
	FirstName string `gorm:"index:idx_full_name"`
	LastName  string `gorm:"index:idx_full_name"`

	// Unique index
	Slug string `gorm:"uniqueIndex"`

	// Ignore field (not stored in database)
	TempData string `gorm:"-"`

	// Read-only field (only read from database)
	CreatedBy string `gorm:"->"`

	// Write-only field (only write to database)
	Secret string `gorm:"<-"`

	// Auto-create and auto-update timestamps
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	// Precision for time fields
	ProcessedAt time.Time `gorm:"precision:6"`

	// Check constraint
	Age int `gorm:"check:age >= 0"`

	// Comment for the column
	Notes string `gorm:"comment:'User notes field'"`

	// Embedded struct
	Address Address `gorm:"embedded;embeddedPrefix:addr_"`
}

// Address is an embedded struct
type Address struct {
	Street  string
	City    string
	Country string
	ZipCode string
}
