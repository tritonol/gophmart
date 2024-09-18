package http

import (
	"encoding/json"
	"net/http"

	"github.com/tritonol/gophmart.git/internal/models/user"
	"github.com/tritonol/gophmart.git/internal/utils/lunh"
)

type withdrawal struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

func (s *Server) GetBalance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctxUserID := ctx.Value("user_id")
	userId, ok := ctxUserID.(user.UserID)

	if !ok {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	balance, err := s.balance.GetBalance(ctx, userId)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	result, err := json.Marshal(balance)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func (s *Server) WriteOff(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctxUserID := ctx.Value("user_id")
	userId, ok := ctxUserID.(user.UserID)
	if !ok {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	var withdraw withdrawal

	if err := json.NewDecoder(r.Body).Decode(&withdraw); err != nil {
		http.Error(w, "can't pasrde body", http.StatusBadRequest)
		return
	}

	balance, err := s.balance.GetBalance(ctx, userId)

	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if balance.Current < withdraw.Sum {
		http.Error(w, "insufficient funds", http.StatusPaymentRequired)
		return
	}

	orderId, err := lunh.Validate(withdraw.Order)
	if err != nil {
		http.Error(w, "wrong number format", http.StatusUnprocessableEntity)
		return
	}

	s.balance.WriteOff(ctx, userId, orderId, -withdraw.Sum)
}

func (s *Server) WithdrawalsHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctxUserID := ctx.Value("user_id")
	userId, ok := ctxUserID.(user.UserID)
	if !ok {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	withdrawals, err :=s.balance.WithdrawalsHistory(ctx, userId)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if len(withdrawals) == 0 {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	result, err := json.Marshal(withdrawals)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}
