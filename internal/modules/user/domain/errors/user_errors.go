package usererrors

import "github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain"

var UserAlreadyExists = domain.FailureError("User already exists")
