package app

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/Desalutar20/lingostruct-server-go/internal/config"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/routes"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
)

func setupRouter(modules *modules, config *config.ApplicationConfig, logger *slog.Logger) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.ClientIPFromRemoteAddr)
	router.Use(
		httprate.LimitBy(100, time.Minute,
			func(r *http.Request) (string, error) {
				return httprate.CanonicalizeIP(middleware.GetClientIP(r.Context())), nil
			},
			httprate.WithLimitHandler(routes.RateLimitHandler)),
	)

	v1 := chi.NewRouter()
	router.Mount("/api/v1", v1)

	middlewares := routes.NewMiddlewares(config, logger, modules.AuthModule.UseCases.AuthenticateCommandHandler)

	modules.AuthModule.ConfigureRoutes(v1, middlewares)
	modules.UserModule.ConfigureRoutes(v1, logger, middlewares)

	return router
}
