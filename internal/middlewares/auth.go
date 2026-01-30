package middlewares

import (
	"context"
	"net/http"
	"time"

	"github.com/Desalutar20/lingostruct-server-go/internal/features/shared"
	"github.com/Desalutar20/lingostruct-server-go/pkg/httputils"
)

func (m *Middlewares) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(m.config.SessionCookieName)

		if err != nil || cookie == nil {
			httputils.ErrorResponse(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		userResponse, err := m.service.Authenticate(r.Context(), cookie.Value, w)
		if err != nil {
			c := &http.Cookie{
				Name:     m.config.SessionCookieName,
				Value:    "",
				Path:     "/",
				Expires:  time.Unix(0, 0),
				HttpOnly: true,
			}

			http.SetCookie(w, c)
			httputils.ErrorResponse(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		ctxWithUser := context.WithValue(r.Context(), shared.UserContextKey, userResponse)
		next.ServeHTTP(w, r.WithContext(ctxWithUser))
	})
}
