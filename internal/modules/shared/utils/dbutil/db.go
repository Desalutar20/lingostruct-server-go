package dbutil

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

func IsDuplicateError(err error) bool {
	e, ok := errors.AsType[*pgconn.PgError](err)

	return ok && e.Code == "23505"
}

func IsDuplicateConstraint(err error, constraint string) bool {
	e, ok := errors.AsType[*pgconn.PgError](err)

	return ok &&
		e.Code == "23505" &&
		e.ConstraintName == constraint
}
