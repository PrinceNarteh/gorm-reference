package handler

import (
	"gorm-reference/internal/service"

	"github.com/gin-gonic/gin"
)

var _ UserHandler = (*userHandler)(nil)

type UserHandler interface {
	Create(*gin.Context)
}

type userHandler struct {
	svc *service.Service
}

func (h *userHandler) Create(c *gin.Context) {
}
