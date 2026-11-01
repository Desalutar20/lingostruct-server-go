package authdto

import (
	userdto "github.com/Desalutar20/lingostruct-server-go/internal/modules/user/dto"
	"github.com/google/uuid"
)

type SignInResult struct {
	UserDto   userdto.UserDto
	SessionId uuid.UUID
}
