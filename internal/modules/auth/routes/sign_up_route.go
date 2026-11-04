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

type SignUpRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (r SignUpRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Email, validation.By(routes.RequiredTrimmed), is.Email, validation.Length(1, valueobjects.EmailMaxLength)),
		validation.Field(&r.Password, validation.By(routes.RequiredTrimmed), validation.Length(userdomain.PasswordMinLength, userdomain.PasswordMaxLength)),
		validation.Field(&r.FirstName, validation.By(routes.RequiredTrimmed), validation.Length(1, userdomain.FirstNameMaxLength)),
		validation.Field(&r.LastName, validation.By(routes.RequiredTrimmed), validation.Length(1, userdomain.LastNameMaxLength)),
	)
}

func (r SignUpRequest) ToCommand() (*authusecases.SignUpCommand, []error) {
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
	firstName, _ := userdomain.NewFirstName(r.FirstName)
	lastName, _ := userdomain.NewLastName(r.LastName)

	return &authusecases.SignUpCommand{
		Email:     email,
		Password:  password,
		FirstName: firstName,
		LastName:  lastName,
	}, nil
}

func (router *AuthRouter) signUp(w http.ResponseWriter, r *http.Request) {

	request, err := routes.DecodeJSON[SignUpRequest](r)
	if routes.WriteError(w, err, router.logger) {
		return
	}

	command, validationErrors := request.ToCommand()
	if routes.WriteErrors(w, validationErrors, router.logger) {
		return
	}

	_, err = router.signUpCommandHandler.Handle(r.Context(), command)
	if routes.WriteError(w, err, router.logger) {
		return
	}

	if err := routes.WriteJSON(w, http.StatusCreated, routes.ResponseDTO[string]{
		Data: "User created successfully",
	}); err != nil {
		router.logger.Error(
			"Failed to write sign-up response",
			slog.Any("error", err),
		)
	}
}
