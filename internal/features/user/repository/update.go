package repository

import (
	"context"
	"time"

	"github.com/Desalutar20/lingostruct-server-go/internal/features/user/model"
)

func (r *Repository) Update(ctx context.Context, user *model.User) error {
	stmt := `
		UPDATE users SET
			first_name = $1,
			last_name = $2,
			email = $3,
			hashed_password = $4,
			avatar_id = $5,
			avatar_url = $6,
			is_verified = $7,
			is_banned = $8,
			updated_at = $9,
			deleted_at = $10
		WHERE id = $11;
	`

	_, err := r.pool.Exec(ctx, stmt, user.FirstName, user.LastName, user.Email, user.HashedPassword, user.AvatarId, user.AvatarUrl, user.IsVerified, user.IsBanned, time.Now().UTC(), user.DeletedAt, user.ID)
	if err != nil {
		return err
	}

	return nil
}
