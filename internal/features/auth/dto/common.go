package dto

import "github.com/Desalutar20/lingostruct-server-go/internal/features/user/dto"

type UserWithSessionId struct {
	SessionId    string
	UserResponse dto.UserResponse
}
