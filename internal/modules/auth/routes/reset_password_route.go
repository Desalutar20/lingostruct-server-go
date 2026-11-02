package authroutes

import (
	"log/slog"
	"net/http"

	authusecases "github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/usecases"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/valueobjects"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/routes"
	userdomain "github.com/Desalutar20/lingostruct-server-go/internal/modules/user/domain"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type ResetPasswordRequest struct {
	Email       string `json:"email"`
	Token       string `json:"token"`
	NewPassword string `json:"newPassword"`
}

func (r ResetPasswordRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Email, validation.Required, is.Email, validation.Length(1, valueobjects.EmailMaxLength)),
		validation.Field(&r.Token, validation.Required),
		validation.Field(&r.NewPassword, validation.Required, validation.Length(userdomain.PasswordMinLength, userdomain.PasswordMaxLength)),
	)
}

func (r ResetPasswordRequest) ToCommand() (*authusecases.ResetPasswordCommand, []error) {
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
	token, _ := valueobjects.NewNonEmptyString(r.Token, "token", "Token", 0, 0, nil)
	newPassword, _ := userdomain.NewPassword(r.NewPassword)

	return &authusecases.ResetPasswordCommand{
		Email:       email,
		Token:       token,
		NewPassword: newPassword,
	}, nil
}

func (router *AuthRouter) resetPassword(w http.ResponseWriter, r *http.Request) {
	request, err := routes.DecodeJSON[ResetPasswordRequest](r)
	if routes.WriteError(w, err, router.logger) {
		return
	}

	command, validationErrors := request.ToCommand()
	if routes.WriteErrors(w, validationErrors, router.logger) {
		return
	}

	_, err = router.resetPasswordCommandHandler.Handle(r.Context(), command)
	if routes.WriteError(w, err, router.logger) {
		return
	}

	if err := routes.WriteJSON(w, http.StatusOK, routes.ResponseDTO[string]{
		Data: "Your password has been reset successfully.",
	}); err != nil {
		router.logger.Error(
			"Failed to write reset-password response",
			slog.Any("error", err),
		)
	}
}
