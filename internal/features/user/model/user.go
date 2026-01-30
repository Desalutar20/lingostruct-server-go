package model

import "time"

type UserRole string

const (
	UserRoleRegular UserRole = "regular"
	UserRoleAdmin   UserRole = "admin"
)

type User struct {
	ID             string     `db:"id"`
	CreatedAt      time.Time  `db:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at"`
	DeletedAt      *time.Time `db:"deleted_at"`
	FirstName      string     `db:"first_name"`
	LastName       string     `db:"last_name"`
	Email          string     `db:"email"`
	HashedPassword string     `db:"hashed_password"`
	Role           UserRole   `db:"role"`
	AvatarId       *string    `db:"avatar_id"`
	AvatarUrl      *string    `db:"avatar_url"`
	IsVerified     bool       `db:"is_verified"`
	IsBanned       bool       `db:"is_banned"`
}
