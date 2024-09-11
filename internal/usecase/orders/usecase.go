package orders

import (
	"context"

	"github.com/tritonol/gophmart.git/internal/models"
)

type orderUsecase struct {
	repo OrderRpository
}

type OrderRpository interface {
	Create(ctx context.Context, model *models.Order) error
	GetUserOrders(ctx context.Context, userId models.UserID)
}

func New(repo OrderRpository) *orderUsecase {
	return &orderUsecase{
		repo: repo,
	}
}

