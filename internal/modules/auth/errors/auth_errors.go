package autherrors

import "github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain"

var ErrInvalidCredentials = domain.FailureError("Invalid credentials")
var ErrUnautorized = domain.UnautorizedError()
