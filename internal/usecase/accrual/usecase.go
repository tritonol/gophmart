package accrual

import (
	"context"
	"fmt"
	"time"

	"github.com/tritonol/gophmart.git/internal/models"
	"github.com/tritonol/gophmart.git/internal/models/order"
)

type accrualUsecase struct {
	orders  orderRepository
	balance balanceRepository
	
	accrual accrualAPI
}

type balanceRepository interface {
	Conduct(ctx context.Context, userID, fromID int64, value float64) error
}

type orderRepository interface {
	GetUnhandledOrders(ctx context.Context) ([]*order.Order, error)
	UpdateStatus(ctx context.Context, orderID int64, status string) error
}

type accrualAPI interface {
	GetAccrual(ctx context.Context, orderID int64) (*models.Accrual, error)
}

func New(orders orderRepository, accrual accrualAPI, balance balanceRepository) *accrualUsecase {
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
		res, err := uc.accrual.GetAccrual(ctx, v.ID)
		if err != nil {
			fmt.Println(err)
			continue 
		}
		if res.Status != string(v.Status) {
			err := uc.orders.UpdateStatus(ctx, v.ID, res.Status)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
		if res.Value > 0 {
			err := uc.balance.Conduct(ctx, int64(v.UserID), v.ID, res.Value)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
	}
}
