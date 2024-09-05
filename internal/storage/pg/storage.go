package pg

import (
	"context"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/tritonol/gophmart.git/internal/storage/migrator"
)

func NewConnection(ctx context.Context, connString string) (*sqlx.DB, error) {
	_, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	db, err := sqlx.Connect("pgx", connString)
	if err != nil {
		return nil, err
	}

	if err := migrator.Migrate(db); err != nil {
		return nil, err
	}

	return db, nil
}
