package order

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"
	models "github.com/tritonol/gophmart.git/internal/models/order"
	"github.com/tritonol/gophmart.git/internal/models/user"
)

type order struct {
	Id         int64   `db:"id"`
	UserId     int64   `db:"user_id"`
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

	var userId int64

	_, err := r.conn.ExecContext(ctx, `
		INSERT INTO orders(id, user_id, status, uploaded_at)
		VALUES ($1, $2, $3, $4) RETURNING user_id
		`,
		order.Id, order.UserId, order.Status, order.UploadedAt,
	)
	fmt.Println(userId, model.UserId)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.UniqueViolation == pgErr.Code {
			existsForUser, err := r.isExistsForUser(ctx, order)
			if err != nil {
				return err
			}
			fmt.Println(existsForUser)
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

func (r *OrderRepo) GetUserOrders(ctx context.Context, userId user.UserID) ([]*models.Order, error) {
	res := make([]order, 0)
	query := `
		SELECT o.*, b.value FROM orders o
		JOIN balance b ON o.id = b.from_id
		WHERE user_id = $1
		ORDER BY uploaded_at
	`
	err := r.conn.SelectContext(ctx, &res, query, userId)
	if err != nil {
		return nil, err
	}

	return toModels(res), nil
}

func (r *OrderRepo) UpdateStatus(ctx context.Context, orderId int64, status string) error {
	_, err := r.conn.ExecContext(ctx, `
		UPDATE orders SET status = $1
		WHERE order_id = $2
		`,
		orderId, status,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *OrderRepo) GetUnhandledOrders(ctx context.Context) ([]*models.Order, error) {
	query := `
		SELECT id, status, user_id FROM orders
		WHERE status != 'INVALID' ADN status != 'PROCESSED'
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
		order.Id, order.UserId,
	).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func toOrder(model *models.Order, uploadedAt time.Time) order {
	return order{
		Id:         model.Id,
		UserId:     int64(model.UserId),
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
		Id:         order.Id,
		Status:     models.OrderStatus(order.Status),
		UserId:     user.UserID(order.UserId),
		Accrual:    order.Accrual,
		UploadedAt: order.UploadedAt,
	}
}
