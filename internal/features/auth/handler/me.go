package handler

import (
	"net/http"

	"github.com/Desalutar20/lingostruct-server-go/internal/features/shared"
	"github.com/Desalutar20/lingostruct-server-go/internal/features/user/dto"
	"github.com/Desalutar20/lingostruct-server-go/pkg/httputils"
)

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(shared.UserContextKey).(*dto.UserResponse)
	if !ok {
		httputils.ErrorResponse(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	httputils.SuccessResponse(w, user, http.StatusOK)
}
