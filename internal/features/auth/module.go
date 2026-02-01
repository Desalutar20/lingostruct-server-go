package auth

import (
	"log/slog"

	"github.com/Desalutar20/lingostruct-server-go/internal/config"
	"github.com/Desalutar20/lingostruct-server-go/internal/features/auth/handler"
	"github.com/Desalutar20/lingostruct-server-go/internal/features/auth/service"
	"github.com/Desalutar20/lingostruct-server-go/internal/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
)

type Module struct {
	handler *handler.Handler
	Service *service.Service
}

func (m *Module) V1(middlewares *middlewares.Middlewares) *chi.Mux {
	router := chi.NewRouter()

	router.Post("/sign-up", m.handler.SignUp)
	router.Post("/verify-account", m.handler.VerifyAccount)
	router.Post("/sign-in", m.handler.SignIn)
	router.Post("/forgot-password", m.handler.ForgotPassword)
	router.Post("/reset-password", m.handler.ResetPassword)

	router.Group(func(r chi.Router) {
		r.Use(middlewares.Authenticate)
		r.Post("/logout", m.handler.Logout)
		r.Get("/me", m.handler.Me)
		r.Patch("/me", m.handler.UpdateProfile)
	})

	return router
}

type Dependencies struct {
	Config      *config.ApplicationConfig
	Repository  service.Repository
	Redis       *redis.Client
	EmailSender service.EmailSender
	Logger      *slog.Logger
}

func New(deps *Dependencies) *Module {
	service := service.New(deps.Config, deps.Repository, deps.Redis, deps.EmailSender)
	handler := handler.New(service, deps.Logger, deps.Config)

	return &Module{handler: handler, Service: service}
}
