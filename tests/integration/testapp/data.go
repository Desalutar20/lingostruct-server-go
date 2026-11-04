package testapp

import (
	authroutes "github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/routes"
	userdomain "github.com/Desalutar20/lingostruct-server-go/internal/modules/user/domain"
	"github.com/brianvoe/gofakeit/v7"
)

func SignUpRequest() authroutes.SignUpRequest {
	return authroutes.SignUpRequest{
		Email:     Email(),
		Password:  Password(),
		FirstName: gofakeit.LetterN(userdomain.FirstNameMaxLength),
		LastName:  gofakeit.LetterN(userdomain.LastNameMaxLength),
	}
}

func SignInRequest() authroutes.SignInRequest {
	return authroutes.SignInRequest{
		Email:    Email(),
		Password: Password(),
	}
}

func ResetPasswordRequest() authroutes.ResetPasswordRequest {
	return authroutes.ResetPasswordRequest{
		Email:       Email(),
		Token:       "token",
		NewPassword: Password(),
	}
}

func Email() string {
	return gofakeit.Email()
}

func Password() string {
	return gofakeit.Password(true, true, true, false, false, userdomain.PasswordMinLength)
}
