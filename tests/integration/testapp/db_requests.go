package testapp

import (
	"errors"
	"testing"

	userrepo "github.com/Desalutar20/lingostruct-server-go/internal/modules/user/repository"
	"github.com/jackc/pgx/v5"
)

const selectUser = `SELECT
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
	FROM users`

func (a *TestApp) GetUserByEmail(
	t *testing.T,
	email string,
) *userrepo.UserDB {
	t.Helper()

	var user userrepo.UserDB

	err := a.pool.QueryRow(
		t.Context(),
		selectUser+" WHERE email = $1",
		email,
	).Scan(
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
		return nil
	}

	if err != nil {
		t.Fatalf("getting user by email: %v", err)
	}

	return &user
}

func (a *TestApp) BanUserByEmail(
	t *testing.T,
	email string,
) {
	t.Helper()

	_, err := a.pool.Exec(t.Context(), "UPDATE users SET is_banned=true WHERE email=$1", email)
	if err != nil {
		t.Fatalf("ban user by email: %v", err)
	}

}
