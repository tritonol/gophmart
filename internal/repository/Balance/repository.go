package balance

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/tritonol/gophmart.git/internal/models/balance"
)

type BalanceRepo struct {
	conn *sqlx.DB
}

func New(ctx context.Context, db *sqlx.DB) *BalanceRepo {
	return &BalanceRepo{
		conn: db,
	}
}

func (r *BalanceRepo) Conduct(ctx context.Context, userID, fromID int64, value float64) error {
	_, err := r.conn.ExecContext(ctx, `
		INSERT INTO balance (user_id, from_id, value)
		VALUES ($1, $2, $3)
		`,
		userID, fromID, value,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *BalanceRepo) GetCurrent(ctx context.Context, userID int64) (float64, error) {
	var sum float64

	err := r.conn.QueryRowContext(
		ctx,
		`SELECT COALESCE(SUM(value), 0) FROM balance WHERE user_id = $1`,
		userID,
	).Scan(&sum)

	if err != nil {
		return 0, err
	}

	return sum, nil
}

func (r *BalanceRepo) GetTotalSpent(ctx context.Context, userID int64) (float64, error) {
	var sum float64

	err := r.conn.QueryRowContext(
		ctx,
		`SELECT COALESCE(ABS(SUM(value)), 0) FROM balance WHERE user_id = $1 AND value < 0`,
		userID,
	).Scan(&sum)

	if err != nil {
		return 0, err
	}

	return sum, nil
}

func (r *BalanceRepo) GetWithdrawals(ctx context.Context, userID int64) ([]*balance.Transaction ,error) {
	withdrawals := make([]*balance.Transaction, 0)
	query := `
		SELECT id, user_id, from_id, value, processed_at 
		FROM balance WHERE user_id = $1 AND value < 0
	`

	err := r.conn.SelectContext(ctx, &withdrawals, query, userID)
	if err != nil {
		return nil, err
	}

	return withdrawals, nil
}
