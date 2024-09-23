package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/tritonol/gophmart.git/internal/models"
)

type Accrual struct {
	client resty.Client
}

func New(address string) *Accrual {
	return &Accrual{
		client: *resty.New().SetBaseURL(address),
	}
}

func (a *Accrual) GetAccrual(ctx context.Context, orderID int64) (*models.Accrual, error) {
	res := &models.Accrual{}

	resp, err := a.client.R().
		SetContext(ctx).
		SetResult(res).
		Get("api/orders/" + strconv.FormatInt(orderID, 10))

	if err != nil {
		return nil, err
	}

	switch resp.StatusCode() {
	case http.StatusNoContent:
		return nil, fmt.Errorf("order not found")
	case http.StatusTooManyRequests:
		return nil, fmt.Errorf("too many requests")
	case http.StatusOK:
		return res, nil
	default:
		return nil, fmt.Errorf("status code not definde")
	}
}
