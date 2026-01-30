package setup

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Desalutar20/lingostruct-server-go/internal/features/user/model"
	"github.com/jackc/pgx/v5"
)

func (a *TestApp) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
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
				WHERE email = $1;
				`

	row, err := a.pool.Query(ctx, query, email)
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

func (a *TestApp) BanUser(ctx context.Context, email string) error {
	stmt := `
	UPDATE users SET is_banned = true
	WHERE email = $1;
				`

	_, err := a.pool.Exec(ctx, stmt, email)
	if err != nil {
		return err
	}

	return nil
}

func (a *TestApp) UnverifyUser(ctx context.Context, email string) error {
	stmt := `
	UPDATE users SET is_verified = false
	WHERE email = $1;
				`

	_, err := a.pool.Exec(ctx, stmt, email)
	if err != nil {
		return err
	}

	return nil
}
