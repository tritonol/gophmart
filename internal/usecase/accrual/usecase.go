package accrual

import (
	"context"
	"fmt"
	"time"

	"github.com/tritonol/gophmart.git/internal/models"
)

type accrualUsecase struct {
	orders  orderRepository
	balance balanceRepository
	
	accrual accrualApi
}

type balanceRepository interface {
	Conduct(ctx context.Context, userId, fromId int64, value float64) error
}

type orderRepository interface {
	GetUnhandledOrders(ctx context.Context) ([]*models.Order, error)
	UpdateStatus(ctx context.Context, orderId int64, status string) error
}

type accrualApi interface {
	GetAccrual(ctx context.Context, orderId int64) (*models.Accrual, error)
}

func New(orders orderRepository, accrual accrualApi, balance balanceRepository) *accrualUsecase {
	return &accrualUsecase{
		orders:  orders,
		accrual: accrual,
		balance: balance,
	}
}

func (uc *accrualUsecase) StartProcessingAccruals(ctx context.Context) {
	go func() {
		delay := time.Second * 5
		for {
			uc.updateOrders(ctx)

			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return
			}
		}
	}()
}

// TODO add logger

func (uc *accrualUsecase) updateOrders(ctx context.Context) {
	unhandledOrders, err := uc.orders.GetUnhandledOrders(ctx)
	if err != nil {
		return
	}

	if len(unhandledOrders) == 0 {
		return
	}

	for _, v := range unhandledOrders {
		res, err := uc.accrual.GetAccrual(ctx, v.Id)
		if err != nil {
			fmt.Println(err)
			continue 
		}
		if res.Status != string(v.Status) {
			err := uc.orders.UpdateStatus(ctx, v.Id, res.Status)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
		if res.Value > 0 {
			err := uc.balance.Conduct(ctx, int64(v.UserId), v.Id, res.Value)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
	}
}
