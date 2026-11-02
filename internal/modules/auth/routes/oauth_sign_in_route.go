package authroutes

import (
	"log/slog"
	"net/http"

	"github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/oauth"
	authusecases "github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/usecases"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/valueobjects"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/routes"
	userdomain "github.com/Desalutar20/lingostruct-server-go/internal/modules/user/domain"
	"github.com/go-chi/chi/v5"
)

func (router *AuthRouter) oauthSignIn(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(router.config.OAauthStateCookieName)
	if routes.WriteError(w, err, router.logger) {
		return
	}

	provider, err := userdomain.NewOAuthProvider(chi.URLParam(r, "provider"))
	if routes.WriteError(w, err, router.logger) {
		return
	}

	code, err := valueobjects.NewNonEmptyString(r.URL.Query().Get("code"), "code", "Code", 0, 0, nil)
	if routes.WriteError(w, err, router.logger) {
		return
	}

	receivedState, err := oauth.NewOAuthState(r.URL.Query().Get("state"))
	if routes.WriteError(w, err, router.logger) {
		return
	}

	expectedState, err := oauth.NewOAuthState(cookie.Value)
	if routes.WriteError(w, err, router.logger) {
		return
	}

	result, err := router.oauthSignInCommandHandler.Handle(r.Context(), &authusecases.OAuthSignInCommand{Provider: provider, Code: code, ReceivedState: receivedState, ExpectedState: expectedState})
	if routes.WriteError(w, err, router.logger) {
		return
	}

	router.SetCookie(w, router.config.SessionCookieName, result.String(), int(router.config.SessionTtl.Seconds()))

	if err := routes.WriteJSON(w, http.StatusOK, routes.ResponseDTO[string]{
		Data: "ok",
	}); err != nil {
		router.logger.Error(
			"Failed to write get-oauth-url response",
			slog.Any("error", err),
		)
	}

}
