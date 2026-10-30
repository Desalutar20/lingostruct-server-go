package valueobjects

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain"
)

type NonEmptyString struct {
	value string
}

func (v *NonEmptyString) Value() string {
	return v.value
}

func NewNonEmptyString(value string, field string, label string, minLength int, maxLength int, additionalCheck func(val string) error) (NonEmptyString, error) {
	trimmed := strings.TrimSpace(value)
	length := utf8.RuneCountInString(trimmed)

	if length == 0 {
		return NonEmptyString{}, domain.ValidationError(field, []string{label + " can not be empty"})
	}

	if minLength != 0 && length < minLength {
		return NonEmptyString{}, domain.ValidationError(field, []string{
			fmt.Sprintf("%s can not be shorter than %d characters", label, minLength),
		})
	}

	if maxLength != 0 && length > maxLength {
		return NonEmptyString{}, domain.ValidationError(field, []string{
			fmt.Sprintf("%s can not be longer than %d characters", label, maxLength),
		})
	}

	if additionalCheck != nil {
		if err := additionalCheck(trimmed); err != nil {
			return NonEmptyString{}, err
		}
	}

	return NonEmptyString{value: trimmed}, nil
}
