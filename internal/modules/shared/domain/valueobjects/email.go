package valueobjects

import (
	"regexp"

	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain"
)

const EmailMaxLength = 100

type Email struct {
	NonEmptyString
}

var emailRegex = regexp.MustCompile(
	`(?i)^[a-z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-z0-9](?:[a-z0-9-]*[a-z0-9])?(?:\.[a-z0-9](?:[a-z0-9-]*[a-z0-9])?)+$`,
)

func NewEmail(value string) (Email, error) {
	return NewNonEmptyStringBase[Email](value, "email", "Email", 0, EmailMaxLength,
		func(value string) error {
			if !emailRegex.MatchString(value) {
				return domain.ValidationError("email", []string{
					"Email format is invalid",
				})
			}

			return nil
		},
		func(value NonEmptyString) Email {
			return Email{NonEmptyString: value}
		})
}
