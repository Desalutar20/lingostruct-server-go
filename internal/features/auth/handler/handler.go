package handler

import (
	"log/slog"

	"github.com/Desalutar20/lingostruct-server-go/internal/config"
	"github.com/Desalutar20/lingostruct-server-go/internal/features/auth/service"
)

type Handler struct {
	service *service.Service
	logger  *slog.Logger
	config  *config.ApplicationConfig
}

func New(service *service.Service, logger *slog.Logger, config *config.ApplicationConfig) *Handler {

	return &Handler{
		service,
		logger,
		config,
	}
}
