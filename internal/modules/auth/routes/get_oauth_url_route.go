package authroutes

import (
	"log/slog"
	"net/http"
	"strings"

	authusecases "github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/usecases"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/valueobjects"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/routes"
	userdomain "github.com/Desalutar20/lingostruct-server-go/internal/modules/user/domain"
	"github.com/go-chi/chi/v5"
)

func (router *Router) getOAuthUrl(w http.ResponseWriter, r *http.Request) {
	provider, err := userdomain.NewOAuthProvider(chi.URLParam(r, "provider"))
	if routes.WriteError(w, err, router.logger) {
		return
	}

	var additionalState *valueobjects.NonEmptyString
	redirectPath := strings.TrimSpace(r.URL.Query().Get("redirectPath"))
	if redirectPath != "" {
		result, err := valueobjects.NewNonEmptyString(redirectPath, "redirectPath", "Redirect path", 0, 0, nil)
		if routes.WriteError(w, err, router.logger) {
			return
		}

		additionalState = &result
	}

	result, err := router.generateOAuthUrlQueryHandler.Handle(r.Context(), &authusecases.GenerateOAuthUrlQuery{Provider: provider, AdditionalState: additionalState})
	if routes.WriteError(w, err, router.logger) {
		return
	}

	router.SetCookie(w, router.config.OAauthStateCookieName, result.State.String(), int(router.config.OAuthStateTtl.Seconds()))

	if err := routes.WriteJSON(w, http.StatusOK, routes.ResponseDTO[string]{
		Data: result.Url.Value(),
	}); err != nil {
		router.logger.Error(
			"Failed to write get-oauth-url response",
			slog.Any("error", err),
		)
	}

}
