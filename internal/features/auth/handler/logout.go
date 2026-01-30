package handler

import (
	"log/slog"
	"net/http"

	"github.com/Desalutar20/lingostruct-server-go/internal/middlewares"
	"github.com/Desalutar20/lingostruct-server-go/pkg/apperror"
	"github.com/Desalutar20/lingostruct-server-go/pkg/httputils"
)

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With(slog.String("request_id", middlewares.GetReqID(r.Context())))

	cookie, err := r.Cookie(h.config.SessionCookieName)
	if err != nil || cookie == nil {
		httputils.ErrorResponse(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if err := h.service.Logout(r.Context(), cookie.Value); err != nil {
		apperror.HandleError(w, err, logger)
		return
	}

	httputils.SuccessResponse(w, "success", http.StatusOK)
}
