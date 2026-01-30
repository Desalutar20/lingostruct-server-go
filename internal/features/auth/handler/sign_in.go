package handler

import (
	"log/slog"
	"net/http"

	"github.com/Desalutar20/lingostruct-server-go/internal/features/auth/dto"
	"github.com/Desalutar20/lingostruct-server-go/internal/middlewares"
	"github.com/Desalutar20/lingostruct-server-go/pkg/apperror"
	"github.com/Desalutar20/lingostruct-server-go/pkg/httputils"
)

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	logger := h.logger.With(slog.String("request_id", middlewares.GetReqID(r.Context())))

	data, err := httputils.ParseBody[dto.SignInRequest](w, r)
	if err != nil {
		return
	}

	response, err := h.service.SignIn(r.Context(), data, w)
	if err != nil {
		apperror.HandleError(w, err, logger)
		return
	}

	httputils.SuccessResponse(w, response.UserResponse, http.StatusOK)
}
