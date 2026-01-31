// Package repository
package repository

import "gorm.io/gorm"

type Repository struct {
	User  UserRepository
	Query QueryRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		User:  &userRepository{db: db},
		Query: &queryRepository{db: db},
	}
}
