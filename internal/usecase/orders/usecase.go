package orders

import (
	"context"

	models "github.com/tritonol/gophmart.git/internal/models/order"
	"github.com/tritonol/gophmart.git/internal/models/user"
)

type orderUsecase struct {
	orders  orderRpository
}

type orderRpository interface {
	Create(ctx context.Context, model *models.Order) error
	GetUserOrders(ctx context.Context, userID user.UserID) ([]*models.Order, error)
}

func New(orders orderRpository) *orderUsecase {
	return &orderUsecase{
		orders:  orders,
	}
}

func (uc *orderUsecase) CreateOrder(ctx context.Context, number int64, userID user.UserID) error {
	order := models.Order {
		ID: number,
		UserID: userID,
		Status: "NEW",
	}
	err := uc.orders.Create(ctx, &order)

	if err != nil {
		return err
	}
	return nil
}

func (uc *orderUsecase) GetUserOrders(ctx context.Context, userID user.UserID) ([]*models.Order, error) {
	res, err := uc.orders.GetUserOrders(ctx, userID)
	if err != nil {
		return nil, err
	}

	return res, nil
}
