package routes

import (
	"errors"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const cursorMaxLength = 200

func RequiredTrimmed(value interface{}) error {
	s, ok := value.(string)
	if !ok || strings.TrimSpace(s) == "" {
		return validation.ErrRequired
	}

	return nil
}

func Trimmed(value interface{}) error {
	s, ok := value.(*string)
	if ok && s != nil && strings.TrimSpace(*s) == "" {
		return validation.ErrRequired
	}

	return nil
}

func ValidateKeysetPagination(value interface{}) error {
	request := value.(KeysetPaginationRequest)

	if request.PrevCursor != nil && request.NextCursor != nil {
		return errors.New("prevCursor and nextCursor cannot be used together")
	}

	return nil
}
