package testapp

import (
	"time"

	authroutes "github.com/Desalutar20/lingostruct-server-go/internal/modules/auth/routes"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/utils/ptrutil"
	userdomain "github.com/Desalutar20/lingostruct-server-go/internal/modules/user/domain"
	userrepo "github.com/Desalutar20/lingostruct-server-go/internal/modules/user/repository"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
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

func usersSeedData() []*userrepo.UserDB {
	users := make([]*userrepo.UserDB, 0, 100)

	for v := range 100 {
		users = append(users, &userrepo.UserDB{
			ID:             uuid.New(),
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
			Email:          Email(),
			FirstName:      ptrutil.To(gofakeit.LetterN(userdomain.FirstNameMaxLength)),
			LastName:       ptrutil.To(gofakeit.LetterN(userdomain.LastNameMaxLength)),
			HashedPassword: nil,
			Role:           userdomain.UserRoleRegular().Value(),
			IsBanned:       v%2 == 0,
			IsVerified:     false,
			AvatarURL:      nil,
			GoogleID:       nil,
			GithubID:       nil,
		})
	}

	return users
}
