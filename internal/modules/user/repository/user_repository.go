package userrepo

import (
	"context"
	"errors"
	"fmt"

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

func (r *UserRepository) GetById(ctx context.Context, id userdomain.UserId) (*userdomain.User, error) {
	return r.getUser(ctx, selectUser+` WHERE id = $1`, id)
}

func (r *UserRepository) GetByEmail(ctx context.Context, email valueobjects.Email) (*userdomain.User, error) {
	return r.getUser(ctx, selectUser+` WHERE email = $1`, email.Value())
}

func (r *UserRepository) Add(ctx context.Context, user *userdomain.User) error {
	userDb := fromDomain(user)

	_, err := r.pool.Exec(ctx, `
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
			return usererrors.UserAlreadyExists
		}

		return domain.InternalError(
			fmt.Sprintf("Failed to add user: %v", err),
		)
	}

	return nil
}

func (r *UserRepository) getUser(ctx context.Context, query string, args ...any) (*userdomain.User, error) {
	var user userDB

	err := r.pool.QueryRow(ctx, query, args...).Scan(
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
