package userroutes

import (
	"log/slog"
	"time"

	"github.com/Desalutar20/lingostruct-server-go/internal/config"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/routes"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
)

type UserRouter struct {
	router *chi.Mux
	logger *slog.Logger
}

func NewV1(router *chi.Mux, logger *slog.Logger, rateLimitConfig *config.RateLimitConfig,
	middlewares *routes.Middlewares) *UserRouter {

	var userRouter = &UserRouter{
		router,
		logger,
	}

	getProfileLimiter := httprate.NewRateLimiter(rateLimitConfig.GetProfile, time.Minute, httprate.WithLimitHandler(routes.RateLimitHandler))

	router.Route("/users", func(r chi.Router) {
		r.With(getProfileLimiter.Handler).With(middlewares.Authenticate).Get("/me", userRouter.GetMe)
	})

	return userRouter
}
