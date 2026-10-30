package domain

type ErrorCode int

const (
	Failure ErrorCode = iota
	NotFound
	Validation
	Conflict
	Unauthorized
	AccessForbidden
	Internal
)

type fieldError struct {
	Field  string
	Errors []string
}

type Error struct {
	message string
	code    ErrorCode
	errors  *fieldError
}

func (e Error) Error() string {
	return e.message
}

func (e Error) Message() string {
	return e.message
}

func (e Error) Code() ErrorCode {
	return e.code
}

func (e Error) Errors() (string, []string, bool) {
	if e.errors == nil {
		return "", nil, false
	}

	return e.errors.Field, e.errors.Errors, true
}

func FailureError(message string) Error {
	return Error{message: message, code: Failure}
}

func NotFoundError(message string) Error {
	return Error{message: message, code: NotFound}
}

func ValidationError(field string, errors []string) Error {
	return Error{
		message: "Validation error",
		code:    Validation,
		errors: &fieldError{
			Field:  field,
			Errors: errors,
		},
	}
}

func ConflictError(message string) Error {
	return Error{message: message, code: Conflict}
}
