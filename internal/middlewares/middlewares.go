package middlewares

import (
	"context"
	"net/http"

	"github.com/Desalutar20/lingostruct-server-go/internal/config"
	"github.com/Desalutar20/lingostruct-server-go/internal/features/user/dto"
)

type Service interface {
	Authenticate(ctx context.Context, sessionId string, w http.ResponseWriter) (*dto.UserResponse, error)
}

type Middlewares struct {
	service Service
	config  *config.ApplicationConfig
}

func New(service Service, config *config.ApplicationConfig) *Middlewares {

	return &Middlewares{
		service,
		config,
	}
}
