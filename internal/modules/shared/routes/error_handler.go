package routes

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain"
)

func WriteError(w http.ResponseWriter, err error, log *slog.Logger) bool {
	if err == nil {
		return false
	}

	status := http.StatusInternalServerError
	dto := ErrorDTO{
		Error: "Internal Server Error",
		Code:  "INTERNAL",
	}

	domainError, ok := errors.AsType[domain.Error](err)
	if ok {
		code := domainError.Code()
		if code != domain.Internal {
			dto.Error = domainError.Message()
			dto.Code = code.String()
		}

		switch code {
		case domain.Failure,
			domain.Validation:
			status = http.StatusBadRequest
		case domain.NotFound:
			status = http.StatusNotFound
		case domain.Conflict:
			status = http.StatusConflict
		case domain.Unauthorized:
			status = http.StatusUnauthorized
		case domain.AccessForbidden:
			status = http.StatusForbidden
		case domain.Internal:
			status = http.StatusInternalServerError
		}

		if field, messages, ok := domainError.Errors(); ok {
			dto.Errors = map[string][]string{
				field: messages,
			}
		}
	}

	if writeErr := WriteJSON(w, status, dto); writeErr != nil {
		log.Error("Failed to write error response", slog.Any("error", writeErr))
	}

	if status == http.StatusInternalServerError {
		log.Error("Internal Server Error", slog.Any("Error", err))
	}

	return true
}

func WriteErrors(w http.ResponseWriter, errs []error, log *slog.Logger) bool {
	if len(errs) == 0 {
		return false
	}

	dto := ErrorDTO{
		Error: "Validation error",
	}

	validationErrors := make(map[string][]string)
	status := http.StatusBadRequest

	for _, err := range errs {
		domainErr, ok := errors.AsType[domain.Error](err)
		if !ok {
			return WriteError(w, err, log)
		}

		switch domainErr.Code() {
		case domain.Validation:
			if field, messages, ok := domainErr.Errors(); ok {
				validationErrors[field] = append(validationErrors[field], messages...)
			}

		default:
			return WriteError(w, err, log)
		}
	}

	dto.Errors = validationErrors

	if writeErr := WriteJSON(w, status, dto); writeErr != nil {
		log.Error("Failed to write error response", slog.Any("error", writeErr))
	}

	return true
}
