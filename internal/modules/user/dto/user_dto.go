package userdto

import "github.com/google/uuid"

type UserDto struct {
	Id        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	FirstName *string   `json:"firstName"`
	LastName  *string   `json:"lastName"`
	Role      string    `json:"role"`
	AvatarUrl *string   `json:"avatarUrl"`
}
