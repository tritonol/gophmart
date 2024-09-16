package user

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
	"github.com/tritonol/gophmart.git/internal/models/user"
)

type UserRepo struct {
	conn *sqlx.DB
}

func New(ctx context.Context, db *sqlx.DB) *UserRepo {
	return &UserRepo{
		conn: db,
	}
}

func (r *UserRepo) Create(ctx context.Context, credentials user.UserCredentials) (user.UserID, error) {
	var id int64
	err := r.conn.QueryRowContext(
		ctx,
		`INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id`,
		credentials.Login, credentials.Password,
	).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.UniqueViolation == pgErr.Code {
			return 0, user.NewUserAlreadyExistsError(credentials.Login, err)
		}
		return 0, err
	}

	return user.UserID(id), nil
}

func (r *UserRepo) CheckByCredentials(ctx context.Context, credentials user.UserCredentials) (user.UserID, error) {
	var id int64

	err := r.conn.QueryRowContext(
		ctx,
		`SELECT id FROM users WHERE login = $1 AND password = $2`,
		credentials.Login, credentials.Password,
	).Scan(&id)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, user.NewUserNotFoundError(credentials.Login, err)
		}
		return 0, err
	}

	return user.UserID(id), nil
}
