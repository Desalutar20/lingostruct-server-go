package adminuserdto

import (
	"time"

	"github.com/google/uuid"
)

type AdminUserDto struct {
	Id         uuid.UUID `json:"id" db:"id"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
	Email      string    `json:"email" db:"email"`
	FirstName  *string   `json:"firstName" db:"first_name"`
	LastName   *string   `json:"lastName" db:"last_name"`
	Role       string    `json:"role" db:"role"`
	AvatarUrl  *string   `json:"avatarUrl" db:"avatar_url"`
	IsBanned   bool      `json:"isBanned" db:"is_banned"`
	IsVerified bool      `json:"isVerified" db:"is_verified"`
	GoogleId   *string   `json:"googleId" db:"google_id"`
	GithubId   *string   `json:"githubId" db:"github_id"`
}
