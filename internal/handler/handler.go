// Package handler
package handler

import "gorm-reference/internal/service"

type Handler struct{}

func NewHandler(s *service.Service) *Handler {
	return &Handler{}
}
