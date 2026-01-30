package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Desalutar20/lingostruct-server-go/internal/features/user/model"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) GetById(ctx context.Context, id string) (*model.User, error) {
	query := `
			SELECT
				id,
				created_at,
				updated_at,
				deleted_at,
				first_name,
				last_name,
				email,
				hashed_password,
				role,
				avatar_id,
				avatar_url,
				is_verified,
				is_banned
			FROM users
			WHERE id = $1;
			`

	row, err := r.pool.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}

	user, err := pgx.CollectOneRow(row, pgx.RowToAddrOfStructByPos[model.User])
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}
