package passwordutil

import (
	userdomain "github.com/Desalutar20/lingostruct-server-go/internal/modules/user/domain"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password *userdomain.Password) (userdomain.HashedPassword, error) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password.Value()),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return userdomain.HashedPassword{}, err
	}

	hashedPassword, err := userdomain.NewHashedPassword(string(hash))
	if err != nil {
		return userdomain.HashedPassword{}, err
	}

	return hashedPassword, nil
}

func VerifyPassword(password *userdomain.Password, hashedPassword *userdomain.HashedPassword) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword.Value()),
		[]byte(password.Value()),
	)

	return err == nil
}
