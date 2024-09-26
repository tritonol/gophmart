package pg

import (
	"context"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/tritonol/gophmart.git/internal/storage/migrator"
)

func NewConnection(ctx context.Context, connString string) (*sqlx.DB, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	db, err := sqlx.ConnectContext(ctx, "pgx", connString)
	if err != nil {
		return nil, err
	}

	if err := migrator.Migrate(db); err != nil {
		return nil, err
	}

	return db, nil
}
