package apperror

import "fmt"

type Kind string

const (
	Validation   Kind = "validation"
	NotFound     Kind = "not_found"
	Conflict     Kind = "conflict"
	Unauthorized Kind = "unauthorized"
	Internal     Kind = "internal"
)

type AppError struct {
	Message string
	Kind    Kind
}

func (e AppError) Error() string {
	return fmt.Sprintf("%s: %s", e.Kind, e.Message)
}

func New(kind Kind, message string) AppError {
	return AppError{
		Message: message,
		Kind:    kind,
	}
}
