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

func (v NonEmptyString) Equal(other NonEmptyString) bool {
	return v.value == other.value
}

func (v NonEmptyString) Value() string {
	return v.value
}

func (a NonEmptyString) Compare(b NonEmptyString) int {
	return strings.Compare(a.value, b.value)
}

func (r NonEmptyString) String() string {
	return r.Value()
}

func NewNonEmptyStringBase[T any](
	value string,
	field string,
	label string,
	minLength int,
	maxLength int,
	additionalCheck func(value string) error,
	factory func(value NonEmptyString) T,
) (T, error) {
	trimmed := strings.TrimSpace(value)
	length := utf8.RuneCountInString(trimmed)

	var zero T

	if length == 0 {
		return zero, domain.ValidationError(field, []string{
			label + " cannot be empty",
		})
	}

	if minLength > 0 && length < minLength {
		return zero, domain.ValidationError(field, []string{
			fmt.Sprintf("%s cannot be shorter than %d characters", label, minLength),
		})
	}

	if maxLength > 0 && length > maxLength {
		return zero, domain.ValidationError(field, []string{
			fmt.Sprintf("%s cannot be longer than %d characters", label, maxLength),
		})
	}

	if additionalCheck != nil {
		if err := additionalCheck(trimmed); err != nil {
			return zero, err
		}
	}

	return factory(NonEmptyString{value: trimmed}), nil
}

func NewNonEmptyString(
	value string,
	field string,
	label string,
	minLength int,
	maxLength int,
	additionalCheck func(string) error,
) (NonEmptyString, error) {
	return NewNonEmptyStringBase(
		value,
		field,
		label,
		minLength,
		maxLength,
		additionalCheck,
		func(value NonEmptyString) NonEmptyString {
			return value
		},
	)
}
