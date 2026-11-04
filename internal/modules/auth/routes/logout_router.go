package authroutes

import (
	"log/slog"
	"net/http"

	authusecases "github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/usecases"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/routes"
)

func (router *Router) logout(w http.ResponseWriter, r *http.Request) {
	session, sessionId, ok := routes.GetSession(r.Context())
	if !ok {
		routes.WriteError(w, domain.UnauthorizedError(), router.logger)
		return
	}

	_, err := router.logoutCommandHandler.Handle(r.Context(), &authusecases.LogoutCommand{
		UserId:    session.UserId,
		SessionId: sessionId,
	})
	if routes.WriteError(w, err, router.logger) {
		return
	}

	router.SetCookie(w, router.config.SessionCookieName, "", 0)

	if err := routes.WriteJSON(w, http.StatusOK, routes.ResponseDTO[string]{
		Data: "success",
	}); err != nil {
		router.logger.Error(
			"Failed to write logout response",
			slog.Any("error", err),
		)
	}

}
