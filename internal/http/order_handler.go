package http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/tritonol/gophmart.git/internal/models/order"
	"github.com/tritonol/gophmart.git/internal/models/user"
	"github.com/tritonol/gophmart.git/internal/utils/lunh"
)

func (s *Server) GetOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctxUserID := ctx.Value("user_id")
	userId, ok := ctxUserID.(user.UserID)

	if !ok {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	orders, err := s.order.GetUserOrders(ctx, userId)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
	}

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

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Can not parse body", http.StatusBadRequest)
		return
	}

	orderId, err := lunh.Validate(string(body))
	if err != nil {
		http.Error(w, "wrong number format", http.StatusUnprocessableEntity)
		return
	}

	ctxUserID := ctx.Value("user_id")
	userId, ok := ctxUserID.(user.UserID)
	if !ok {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	err = s.order.CreateOrder(ctx, orderId, userId)
	if err != nil {
		if errors.As(err, order.ErrAlreadyExists) {
			http.Error(w, "", http.StatusConflict)
			return
		}

		if errors.As(err, order.ErrCreatedByAnotherUser) {
			http.Error(w, "", http.StatusConflict)
			return
		}

		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
