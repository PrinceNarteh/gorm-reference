// Package service
package service

import "gorm-reference/internal/repository"

type Service struct {
	User UserService
}

func NewService(r *repository.Repository) *Service {
	return &Service{
		User: userService{repo: r},
	}
}
