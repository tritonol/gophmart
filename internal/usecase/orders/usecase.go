package orders

import (
	"context"

	"github.com/tritonol/gophmart.git/internal/models"
)

type orderUsecase struct {
	orders  orderRpository
}

type orderRpository interface {
	Create(ctx context.Context, model *models.Order) error
	GetUserOrders(ctx context.Context, userId models.UserID) ([]*models.Order, error)
}

func New(orders orderRpository) *orderUsecase {
	return &orderUsecase{
		orders:  orders,
	}
}

func (uc *orderUsecase) CreateOrder(ctx context.Context, order models.Order) error {
	err := uc.orders.Create(ctx, &order)

	if err != nil {
		return err
	}
	return nil
}

func (uc *orderUsecase) GetUserOrders(ctx context.Context, userId models.UserID) ([]*models.Order, error) {
	res, err := uc.orders.GetUserOrders(ctx, userId)
	if err != nil {
		return nil, err
	}

	return res, nil
}
