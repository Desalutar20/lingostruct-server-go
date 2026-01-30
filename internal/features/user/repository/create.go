package repository

import (
	"context"

	"github.com/Desalutar20/lingostruct-server-go/internal/features/user/model"
)

func (r *Repository) Create(ctx context.Context, user *model.User) (string, error) {
	stmt := `INSERT INTO users (first_name, last_name, email, hashed_password) VALUES($1, $2, $3, $4) RETURNING id;`

	var id string
	err := r.pool.QueryRow(context.Background(), stmt, user.FirstName, user.LastName, user.Email, user.HashedPassword).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}
