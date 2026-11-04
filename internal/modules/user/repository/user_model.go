package userrepo

import (
	"time"

	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/utils/ptrutil"
	userdomain "github.com/Desalutar20/lingostruct-server-go/internal/modules/user/domain"
	"github.com/google/uuid"
)

type UserDB struct {
	ID             uuid.UUID `db:"id"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
	FirstName      *string   `db:"first_name"`
	LastName       *string   `db:"last_name"`
	Email          string    `db:"email"`
	HashedPassword *string   `db:"hashed_password"`
	Role           string    `db:"role"`
	IsBanned       bool      `db:"is_banned"`
	IsVerified     bool      `db:"is_verified"`
	AvatarURL      *string   `db:"avatar_url"`
	GoogleID       *string   `db:"google_id"`
	GithubID       *string   `db:"github_id"`
}

func (u UserDB) ToDomain() *userdomain.User {
	return ptrutil.To(userdomain.RestoreUser(
		u.ID,
		u.CreatedAt,
		u.UpdatedAt,
		u.Email,
		u.HashedPassword,
		u.FirstName,
		u.LastName,
		u.Role,
		u.IsBanned,
		u.IsVerified,
		u.AvatarURL,
		u.GoogleID,
		u.GithubID,
	))
}

func fromDomain(user *userdomain.User) UserDB {
	return UserDB{
		uuid.UUID(user.Id()),
		user.CreatedAt(),
		user.UpdatedAt(),
		ptrutil.ValuePtr(user.FirstName()),
		ptrutil.ValuePtr(user.LastName()),
		user.Email().Value(),
		ptrutil.ValuePtr(user.HashedPassword()),
		user.Role().Value(),
		user.IsBanned(),
		user.IsVerified(),
		ptrutil.ValuePtr(user.AvatarUrl()),
		ptrutil.ValuePtr(user.GoogleId()),
		ptrutil.ValuePtr(user.GithubId()),
	}
}
