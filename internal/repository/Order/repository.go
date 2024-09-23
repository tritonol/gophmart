package order

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
	models "github.com/tritonol/gophmart.git/internal/models/order"
	"github.com/tritonol/gophmart.git/internal/models/user"
)

type order struct {
	ID         int64   `db:"id"`
	UserID     int64   `db:"user_id"`
	Status     string  `db:"status"`
	Accrual    float64 `db:"value"`
	UploadedAt string  `db:"uploaded_at"`
}

type OrderRepo struct {
	conn *sqlx.DB
}

func New(ctx context.Context, db *sqlx.DB) *OrderRepo {
	return &OrderRepo{
		conn: db,
	}
}

func (r *OrderRepo) Create(ctx context.Context, model *models.Order) error {
	order := toOrder(model, time.Now())

	_, err := r.conn.ExecContext(ctx, `
		INSERT INTO orders(id, user_id, status, uploaded_at)
		VALUES ($1, $2, $3, $4) RETURNING user_id
		`,
		order.ID, order.UserID, order.Status, order.UploadedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.UniqueViolation == pgErr.Code {
			existsForUser, err := r.isExistsForUser(ctx, order)
			if err != nil {
				return err
			}
			if existsForUser {
				return models.ErrAlreadyExists
			} else {
				return models.ErrCreatedByAnotherUser
			}
		}

		return err
	}
	return nil
}

func (r *OrderRepo) GetUserOrders(ctx context.Context, userID user.UserID) ([]*models.Order, error) {
	res := make([]order, 0)
	query := `
		SELECT o.*, COALESCE(b.value, 0) AS value FROM orders o
		LEFT JOIN balance b ON o.id = b.from_id
		WHERE o.user_id = $1
		ORDER BY uploaded_at
	`
	err := r.conn.SelectContext(ctx, &res, query, userID)

	if err != nil {
		return nil, err
	}

	return toModels(res), nil
}

func (r *OrderRepo) UpdateStatus(ctx context.Context, orderID int64, status string) error {
	_, err := r.conn.ExecContext(ctx, `
		UPDATE orders SET status = $2
		WHERE id = $1
		`,
		orderID, status,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *OrderRepo) GetUnhandledOrders(ctx context.Context) ([]*models.Order, error) {
	query := `
		SELECT id, status, user_id FROM orders
		WHERE status != 'INVALID' AND status != 'PROCESSED'
	`
	orders := make([]order, 0)

	err := r.conn.SelectContext(ctx, &orders, query)
	if err != nil {
		return nil, err
	}

	return toModels(orders), nil
}

func (r *OrderRepo) isExistsForUser(ctx context.Context, order order) (bool, error) {
	var exists bool

	err := r.conn.QueryRowContext(
		ctx,
		`SELECT EXISTS(SELECT 1 FROM orders WHERE id = $1 AND user_id = $2)`,
		order.ID, order.UserID,
	).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func toOrder(model *models.Order, uploadedAt time.Time) order {
	return order{
		ID:         model.ID,
		UserID:     int64(model.UserID),
		Status:     string(model.Status),
		UploadedAt: uploadedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func toModels(orders []order) []*models.Order {
	res := make([]*models.Order, len(orders))
	for i := 0; i < len(orders); i++ {
		res[i] = toModel(orders[i])
	}

	return res
}

func toModel(order order) *models.Order {
	return &models.Order{
		ID:         order.ID,
		Status:     models.OrderStatus(order.Status),
		UserID:     user.UserID(order.UserID),
		Accrual:    order.Accrual,
		UploadedAt: order.UploadedAt,
	}
}
