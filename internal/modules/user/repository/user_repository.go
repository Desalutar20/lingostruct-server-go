package userrepo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Desalutar20/lingostruct-server-go/internal/database/postgres"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/domain/valueobjects"
	"github.com/Desalutar20/lingostruct-server-go/internal/modules/shared/utils/dbutil"
	userdomain "github.com/Desalutar20/lingostruct-server-go/internal/modules/user/domain"
	usererrors "github.com/Desalutar20/lingostruct-server-go/internal/modules/user/domain/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

const selectUser = `
	SELECT
		id,
		created_at,
		updated_at,
		first_name,
		last_name,
		email,
		hashed_password,
		role,
		is_banned,
		is_verified,
		avatar_url,
		google_id,
		github_id
	FROM users
`

func (r *UserRepository) GetById(ctx context.Context, id userdomain.UserId, tx pgx.Tx) (*userdomain.User, error) {
	return r.getUser(ctx, tx, selectUser+` WHERE id = $1`, id)
}

func (r *UserRepository) GetByEmail(ctx context.Context, email valueobjects.Email, tx pgx.Tx) (*userdomain.User, error) {
	return r.getUser(ctx, tx, selectUser+` WHERE email = $1`, email.Value())
}

func (r *UserRepository) Add(ctx context.Context, user *userdomain.User, tx pgx.Tx) error {
	var q postgres.Querier = r.pool
	if tx != nil {
		q = tx
	}

	userDb := fromDomain(user)

	_, err := q.Exec(ctx, `
		INSERT INTO users (
			id,
			created_at,
			updated_at,
			first_name,
			last_name,
			email,
			hashed_password,
			role,
			is_banned,
			is_verified,
			avatar_url,
			google_id,
			github_id
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7,
			$8, $9, $10, $11, $12, $13
		)
	`,
		userDb.ID,
		userDb.CreatedAt,
		userDb.UpdatedAt,
		userDb.FirstName,
		userDb.LastName,
		userDb.Email,
		userDb.HashedPassword,
		userDb.Role,
		userDb.IsBanned,
		userDb.IsVerified,
		userDb.AvatarURL,
		userDb.GoogleID,
		userDb.GithubID,
	)

	if err != nil {
		if dbutil.IsDuplicateError(err) {
			return usererrors.ErrUserAlreadyExists
		}

		return domain.InternalError(
			fmt.Sprintf("Failed to add user: %v", err),
		)
	}

	return nil
}

func (r *UserRepository) Update(ctx context.Context, user *userdomain.User, tx pgx.Tx) error {
	var q postgres.Querier = r.pool
	if tx != nil {
		q = tx
	}
	userDb := fromDomain(user)

	_, err := q.Exec(ctx, `
		UPDATE users SET
			updated_at = $1,
			first_name = $2,
			last_name = $3,
			email = $4,
			hashed_password = $5,
			role = $6,
			is_banned = $7,
			is_verified = $8,
			avatar_url = $9,
			google_id = $10,
			github_id = $11
		WHERE id = $12
	`,
		time.Now(),
		userDb.FirstName,
		userDb.LastName,
		userDb.Email,
		userDb.HashedPassword,
		userDb.Role,
		userDb.IsBanned,
		userDb.IsVerified,
		userDb.AvatarURL,
		userDb.GoogleID,
		userDb.GithubID,
		user.Id(),
	)

	if err != nil {
		return domain.InternalError(
			fmt.Sprintf("Failed to update user: %v", err),
		)
	}

	return nil
}

func (r *UserRepository) Remove(ctx context.Context, user *userdomain.User, tx pgx.Tx) error {
	var q postgres.Querier = r.pool
	if tx != nil {
		q = tx
	}

	_, err := q.Exec(ctx, `
		DELETE FROM users WHERE id = $1
	`,
		user.Id(),
	)

	if err != nil {
		return domain.InternalError(
			fmt.Sprintf("Failed to remove user: %v", err),
		)
	}

	return nil
}

func (r *UserRepository) getUser(ctx context.Context, tx pgx.Tx, query string, args ...any) (*userdomain.User, error) {
	var q postgres.Querier = r.pool
	if tx != nil {
		q = tx
	}
	var user UserDB

	err := q.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.HashedPassword,
		&user.Role,
		&user.IsBanned,
		&user.IsVerified,
		&user.AvatarURL,
		&user.GoogleID,
		&user.GithubID,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, domain.InternalError(
			fmt.Sprintf("Failed to get user: %v", err),
		)
	}

	return user.ToDomain(), nil
}
