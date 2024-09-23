package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/tritonol/gophmart.git/internal/models/order"
	"github.com/tritonol/gophmart.git/internal/models/user"
	"github.com/tritonol/gophmart.git/internal/utils/lunh"
)

type respOrder struct {
	Number     string  `json:"number"`
	Accrual    float64 `json:"accrual"`
	Status     string  `json:"status"`
	UploadedAt string  `json:"uploaded_at"`
}

func (s *Server) GetOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctxUserID := ctx.Value(keyUserID)
	userID, ok := ctxUserID.(user.UserID)

	if !ok {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	rawOrders, err := s.order.GetUserOrders(ctx, userID)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
	}

	orders := toRespOrders(rawOrders)

	if len(orders) == 0 {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	result, err := json.Marshal(orders)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func (s *Server) CreateOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctxUserID := ctx.Value(keyUserID)
	userID, ok := ctxUserID.(user.UserID)
	if !ok {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Can not parse body", http.StatusBadRequest)
		return
	}

	orderID, err := lunh.Validate(string(body))
	if err != nil {
		http.Error(w, "wrong number format", http.StatusUnprocessableEntity)
		return
	}

	err = s.order.CreateOrder(ctx, orderID, userID)
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, order.ErrAlreadyExists) {
			http.Error(w, "", http.StatusOK)
			return
		}

		if errors.Is(err, order.ErrCreatedByAnotherUser) {
			http.Error(w, "", http.StatusConflict)
			return
		}

		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func toRespOrders(ords []*order.Order) []*respOrder {
	ordsLen := len(ords)
	res := make([]*respOrder, ordsLen)

	for i := 0; i < ordsLen; i++ {
		res[i] = toRespOrder(ords[i])
	}

	return res
}

func toRespOrder(ord *order.Order) *respOrder {
	return &respOrder{
		Number: strconv.FormatInt(ord.ID, 10),
		Accrual: ord.Accrual,
		Status: string(ord.Status),
		UploadedAt: ord.UploadedAt,
	}
}