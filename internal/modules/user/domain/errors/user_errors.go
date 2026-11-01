package usererrors

import "github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain"

var ErrUserAlreadyExists = domain.FailureError("User already exists")
