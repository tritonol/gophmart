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

func (r *BalanceRepo) Conduct(ctx context.Context, userId, fromId int64, value float64) error {
	_, err := r.conn.ExecContext(ctx, `
		INSERT INTO balance (user_id, from_id, value)
		VALUES ($1, $2, $3)
		`,
		userId, fromId, value,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *BalanceRepo) GetCurrent(ctx context.Context, userId int64) (float64, error) {
	var sum float64

	err := r.conn.QueryRowContext(
		ctx,
		`SELECT COALESCE(SUM(value), 0) FROM balance WHERE user_id = $1`,
		userId,
	).Scan(&sum)

	if err != nil {
		return 0, err
	}

	return sum, nil
}

func (r *BalanceRepo) GetTotalSpent(ctx context.Context, userId int64) (float64, error) {
	var sum float64

	err := r.conn.QueryRowContext(
		ctx,
		`SELECT COALESCE(ABS(SUM(value)), 0) FROM balance WHERE user_id = $1 AND value < 0`,
		userId,
	).Scan(&sum)

	if err != nil {
		return 0, err
	}

	return sum, nil
}

func (r *BalanceRepo) GetWithdrawals(ctx context.Context, userId int64) ([]*balance.Transaction ,error) {
	withdrawals := make([]*balance.Transaction, 0)
	query := `
		SELECT id, user_id, from_id, value, processed_at 
		FROM balance WHERE user_id = $1
		WHERE value < 0
	`

	err := r.conn.SelectContext(ctx, &withdrawals, query, userId)
	if err != nil {
		return nil, err
	}

	return withdrawals, nil
}
