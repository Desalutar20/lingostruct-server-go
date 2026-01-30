package dto

import (
	"time"

	"github.com/Desalutar20/lingostruct-server-go/internal/features/user/model"
)

type UserResponse struct {
	ID        string         `json:"id"`
	DeletedAt *time.Time     `json:"deleted_at"`
	FirstName string         `json:"first_name"`
	LastName  string         `json:"last_name"`
	Email     string         `json:"email"`
	Role      model.UserRole `json:"role"`
	AvatarId  *string        `json:"avatar_id"`
	AvatarUrl *string        `json:"avatar_url"`
}
