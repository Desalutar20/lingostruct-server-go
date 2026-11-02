package autherrors

import "github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain"

var (
	ErrInvalidCredentials = domain.FailureError("Invalid credentials")
	ErrInvalidToken       = domain.FailureError("Token is invalid")
	ErrUnauthorized       = domain.UnauthorizedError()
)
