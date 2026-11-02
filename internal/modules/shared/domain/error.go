package domain

type ErrorCode int

func (c ErrorCode) String() string {
	switch c {
	case Failure:
		return "BAD_REQUEST"
	case NotFound:
		return "NOT_FOUND"
	case Validation:
		return "VALIDATION"
	case Conflict:
		return "CONFLICT"
	case Unauthorized:
		return "UNAUTHORIZED"
	case AccessForbidden:
		return "ACCESS_FORBIDDEN"
	default:
		return "INTERNAL"
	}
}

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

func UnauthorizedError() Error {
	return Error{message: "Unauthorized", code: Unauthorized}
}
func AccessForbiddenError() Error {
	return Error{message: "Access forbidden", code: AccessForbidden}
}

func InternalError(message string) Error {
	return Error{message: message, code: Internal}
}
