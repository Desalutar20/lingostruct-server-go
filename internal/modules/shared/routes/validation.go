package routes

import (
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func RequiredTrimmed(value interface{}) error {
	s, ok := value.(string)
	if !ok || strings.TrimSpace(s) == "" {
		return validation.ErrRequired
	}

	return nil
}
