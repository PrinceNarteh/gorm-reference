// Package repository
package repository

import "gorm.io/gorm"

type Repository struct {
	User UserRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		User: &userRepository{db: db},
	}
}
