package balance

import (
	"context"

	"github.com/tritonol/gophmart.git/internal/models/balance"
	"github.com/tritonol/gophmart.git/internal/models/user"
)

type balanceUsecase struct {
	balance balanceRepository
}

type balanceRepository interface {
	Conduct(ctx context.Context, userID, fromID int64, value float64) error
	GetCurrent(ctx context.Context, userID int64) (float64, error)
	GetTotalSpent(ctx context.Context, userID int64) (float64, error)
	GetWithdrawals(ctx context.Context, userID int64) ([]*balance.Transaction, error)
}

func New(balance balanceRepository) *balanceUsecase {
	return &balanceUsecase{
		balance: balance,
	}
}

func (uc *balanceUsecase) GetBalance(ctx context.Context, userID user.UserID) (*balance.Balance, error) {
	current, err := uc.balance.GetCurrent(ctx, int64(userID))
	if err != nil {
		return nil, err
	}

	withdrawn, err := uc.balance.GetTotalSpent(ctx, int64(userID))
	if err != nil {
		return nil, err
	}

	return &balance.Balance{
		Current:   current,
		Withdrawn: withdrawn,
	}, nil
}

func (uc *balanceUsecase) WriteOff(ctx context.Context, userID user.UserID, orderNum int64, value float64) error {
	currentAmount, err := uc.balance.GetCurrent(ctx, int64(userID))
	if err != nil {
		return err
	}

	if currentAmount < value {
		return balance.ErrInsufficientFunds
	}

	err = uc.balance.Conduct(ctx, int64(userID), orderNum, value)
	if err != nil {
		return err
	}

	return nil
}

func (uc *balanceUsecase) WithdrawalsHistory(ctx context.Context, userID user.UserID) ([]*balance.Transaction, error) {
	history, err := uc.balance.GetWithdrawals(ctx, int64(userID))
	if err != nil {
		return nil, err
	}

	return history, nil
}