package authroutes

import (
	"log/slog"
	"net/http"

	authusecases "github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/usecases"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/valueobjects"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/routes"
	userdomain "github.com/Desalutar20/lingostruct-server-go/internal/modules/user/domain"
	userdto "github.com/Desalutar20/lingostruct-server-go/internal/modules/user/dto"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type SignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r SignInRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Email, validation.By(routes.RequiredTrimmed), is.Email, validation.Length(1, valueobjects.EmailMaxLength)),
		validation.Field(&r.Password, validation.By(routes.RequiredTrimmed), validation.Length(userdomain.PasswordMinLength, userdomain.PasswordMaxLength)),
	)
}

func (r SignInRequest) ToCommand() (*authusecases.SignInCommand, []error) {
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
	password, _ := userdomain.NewPassword(r.Password)

	return &authusecases.SignInCommand{
		Email:    email,
		Password: password,
	}, nil
}

func (router *AuthRouter) signIn(w http.ResponseWriter, r *http.Request) {
	request, err := routes.DecodeJSON[SignInRequest](r)
	if routes.WriteError(w, err, router.logger) {
		return
	}

	command, validationErrors := request.ToCommand()
	if routes.WriteErrors(w, validationErrors, router.logger) {
		return
	}

	result, err := router.signInCommandHandler.Handle(r.Context(), command)
	if routes.WriteError(w, err, router.logger) {
		return
	}

	router.SetCookie(w, router.config.SessionCookieName, result.SessionId.String(), int(router.config.SessionTtl.Seconds()))

	if err := routes.WriteJSON(w, http.StatusOK, routes.ResponseDTO[userdto.UserDto]{
		Data: result.UserDto,
	}); err != nil {
		router.logger.Error(
			"Failed to write sign-in response",
			slog.Any("error", err),
		)
	}
}
