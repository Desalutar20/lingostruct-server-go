package apperror

import (
	"log/slog"
	"net/http"

	"github.com/Desalutar20/lingostruct-server-go/pkg/httputils"
)

func HandleError(w http.ResponseWriter, err error, log *slog.Logger) {
	apiErr, ok := err.(AppError)
	if ok {
		code := mapKindToStatus(apiErr.Kind)
		httputils.ErrorResponse(w, apiErr.Message, code)
		return
	}

	log.Error("unexpected error", slog.Any("err", err))
	httputils.ErrorResponse(w, "something went wrong", http.StatusInternalServerError)
}

func mapKindToStatus(kind Kind) int {
	switch kind {
	case Validation:
		return http.StatusBadRequest
	case Conflict:
		return http.StatusBadRequest
	case NotFound:
		return http.StatusNotFound
	case Unauthorized:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}
