package authroutes

import (
	"log/slog"
	"net/http"

	authusecases "github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/usecases"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/valueobjects"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/routes"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

func (r ForgotPasswordRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Email, validation.Required, is.Email, validation.Length(1, valueobjects.EmailMaxLength)),
	)
}

func (r ForgotPasswordRequest) ToCommand() (*authusecases.ForgotPasswordCommand, []error) {
	err := r.Validate()
	if err != nil {
		flattenErrors := validation.Flatten(err)
		errStrings := make([]error, 0, len(flattenErrors))

		for _, fe := range flattenErrors {
			errStrings = append(errStrings, domain.ValidationError(fe.Field(), []string{fe.Message}))
		}

		return nil, errStrings

	}

	email, _ := valueobjects.NewEmail(r.Email)

	return &authusecases.ForgotPasswordCommand{
		Email: email,
	}, nil
}

func (router *AuthRouter) forgotPassword(w http.ResponseWriter, r *http.Request) {
	request, err := routes.DecodeJSON[ForgotPasswordRequest](r)
	if routes.WriteError(w, err, router.logger) {
		return
	}

	command, validationErrors := request.ToCommand()
	if routes.WriteErrors(w, validationErrors, router.logger) {
		return
	}

	_, err = router.forgotPasswordCommandHandler.Handle(r.Context(), command)
	if routes.WriteError(w, err, router.logger) {
		return
	}

	if err := routes.WriteJSON(w, http.StatusOK, routes.ResponseDTO[string]{
		Data: "If the email is registered, password reset instructions have been sent",
	}); err != nil {
		router.logger.Error(
			"Failed to write forgot-password response",
			slog.Any("error", err),
		)
	}
}
