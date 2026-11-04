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

type VerifyAccountRequest struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

func (r VerifyAccountRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Email, validation.By(routes.RequiredTrimmed), is.Email, validation.Length(1, valueobjects.EmailMaxLength)),
		validation.Field(&r.Token, validation.By(routes.RequiredTrimmed)),
	)
}

func (r VerifyAccountRequest) ToCommand() (*authusecases.VerifyAccountCommand, []error) {
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

	return &authusecases.VerifyAccountCommand{
		Email: email,
		Token: token,
	}, nil
}

func (router *AuthRouter) verifyAccount(w http.ResponseWriter, r *http.Request) {
	request, err := routes.DecodeJSON[VerifyAccountRequest](r)
	if routes.WriteError(w, err, router.logger) {
		return
	}

	command, validationErrors := request.ToCommand()
	if routes.WriteErrors(w, validationErrors, router.logger) {
		return
	}

	_, err = router.verifyAccountCommandHandler.Handle(r.Context(), command)
	if routes.WriteError(w, err, router.logger) {
		return
	}

	if err := routes.WriteJSON(w, http.StatusOK, routes.ResponseDTO[string]{
		Data: "Your account has been successfully verified. You can now log in.",
	}); err != nil {
		router.logger.Error(
			"Failed to write verify-account response",
			slog.Any("error", err),
		)
	}
}
