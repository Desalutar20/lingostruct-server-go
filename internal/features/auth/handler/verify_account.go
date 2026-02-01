package handler

import (
	"log/slog"
	"net/http"

	"github.com/Desalutar20/lingostruct-server-go/internal/features/auth/dto"
	"github.com/Desalutar20/lingostruct-server-go/internal/middlewares"
	"github.com/Desalutar20/lingostruct-server-go/pkg/apperror"
	"github.com/Desalutar20/lingostruct-server-go/pkg/httputils"
)

func (h *Handler) VerifyAccount(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With(slog.String("request_id", middlewares.GetReqID(r.Context())))

	data, err := httputils.ParseData[dto.VerifyAccountRequest](w, r.Body)
	if err != nil {
		return
	}

	if err := h.service.VerifyAccount(r.Context(), data); err != nil {
		apperror.HandleError(w, err, logger)
		return
	}

	httputils.SuccessResponse(w, "Your email has been verified! You can now log in.", http.StatusOK)
}
