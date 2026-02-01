package handler

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Desalutar20/lingostruct-server-go/internal/features/auth/constants"
	"github.com/Desalutar20/lingostruct-server-go/internal/features/auth/dto"
	"github.com/Desalutar20/lingostruct-server-go/internal/features/shared"
	userDto "github.com/Desalutar20/lingostruct-server-go/internal/features/user/dto"
	"github.com/Desalutar20/lingostruct-server-go/internal/middlewares"
	"github.com/Desalutar20/lingostruct-server-go/pkg/apperror"
	"github.com/Desalutar20/lingostruct-server-go/pkg/httputils"
)

func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(shared.UserContextKey).(*userDto.UserResponse)
	if !ok {
		httputils.ErrorResponse(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	extraBytesForFields := 2048
	r.Body = http.MaxBytesReader(w, r.Body, int64(constants.MaxFileSize+extraBytesForFields))

	logger := h.logger.With(slog.String("request_id", middlewares.GetReqID(r.Context())))

	err := r.ParseMultipartForm(constants.MaxFileSize)
	if err != nil {
		if errors.Is(err, http.ErrNotMultipart) {
			httputils.ErrorResponse(w, "request must be multipart/form-data", http.StatusUnsupportedMediaType)
		} else {
			apperror.HandleError(w, err, logger)
		}

		return
	}
	defer r.MultipartForm.RemoveAll()

	data := &dto.UpdateProfileRequest{}

	file, handler, err := r.FormFile("image")
	if err != nil {
		if err.Error() != "http: no such file" {
			apperror.HandleError(w, err, logger)
			return
		}
	}

	if handler != nil {
		defer file.Close()
		fmt.Println(file)
		data.Image = file
		fmt.Println(file)

		mimeType := handler.Header.Get("Content-Type")
		if mimeType != "image/png" && mimeType != "image/webp" && mimeType != "image/jpeg" {
			httputils.ErrorResponse(w, "invalid mime type", http.StatusBadRequest)
			return
		}

	}

	for key, value := range r.PostForm {
		if key == "firstName" && len(value) > 0 {
			data.FirstName = value[0]
		}

		if key == "lastName" && len(value) > 0 {
			data.LastName = value[0]
		}
	}

	fmt.Printf("data: %v\n", data)

	httputils.SuccessResponse(w, "Check your email to verify your account", http.StatusOK)
}
