package handlers

import "REST_API_Songs/internal/service"

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}
