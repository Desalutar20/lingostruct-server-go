package routes

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/Desalutar20/lingostruct-server-go/internal/config"
	authtypes "github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/types"
	authusecases "github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/usecases"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/cqrs"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain"
	"github.com/google/uuid"
)

type sessionContextKey struct{}
type sessionIdContextKey struct{}

type Middlewares struct {
	config *config.ApplicationConfig
	logger *slog.Logger

	authenticateCommandHandler cqrs.CommandHandler[*authusecases.AuthenticateCommand, *authtypes.Session]
}

func (m *Middlewares) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(m.config.SessionCookieName)
		if WriteError(w, err, m.logger) {
			return
		}

		sessionID, err := uuid.Parse(cookie.Value)
		if err != nil {
			m.logger.Warn("invalid session ID")
			WriteError(w, domain.UnauthorizedError(), m.logger)
			return
		}

		session, err := m.authenticateCommandHandler.Handle(
			r.Context(),
			&authusecases.AuthenticateCommand{
				SessionId: sessionID,
			},
		)
		if WriteError(w, err, m.logger) {
			return
		}

		ctx := context.WithValue(r.Context(), sessionContextKey{}, session)
		ctx = context.WithValue(ctx, sessionIdContextKey{}, sessionID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetSession(ctx context.Context) (*authtypes.Session, uuid.UUID, bool) {
	session, sessionOK := ctx.Value(sessionContextKey{}).(*authtypes.Session)
	sessionID, sessionIDOK := ctx.Value(sessionIdContextKey{}).(uuid.UUID)

	return session, sessionID, sessionOK && sessionIDOK
}

func NewMiddlewares(
	config *config.ApplicationConfig,
	logger *slog.Logger,
	authenticateCommandHandler cqrs.CommandHandler[*authusecases.AuthenticateCommand, *authtypes.Session],
) *Middlewares {
	return &Middlewares{
		config, logger, authenticateCommandHandler,
	}
}
