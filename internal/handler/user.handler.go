package handler

import "gorm-reference/internal/service"

var _ UserHandler = (*userHandler)(nil)

type UserHandler interface{}

type userHandler struct {
	svc *service.Service
}

func (h *userHandler) Create() {}
