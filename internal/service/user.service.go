package service

import (
	"context"

	"gorm-reference/internal/models"
	"gorm-reference/internal/repository"
)

var _ UserService = (*userService)(nil)

type UserService interface {
	Create(ctx context.Context, user *models.User) error
}

type userService struct {
	repo *repository.Repository
}

func (s *userService) Create(ctx context.Context, user *models.User) error {
	return s.repo.User.Create(ctx, user)
}
