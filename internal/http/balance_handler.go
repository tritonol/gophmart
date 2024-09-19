package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/tritonol/gophmart.git/internal/models/balance"
	"github.com/tritonol/gophmart.git/internal/models/user"
	"github.com/tritonol/gophmart.git/internal/utils/lunh"
)

type withdrawal struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
	Processed_at string `json:"processed_at"`
}

func (s *Server) GetBalance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctxUserID := ctx.Value(keyUserId)
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

	ctxUserID := ctx.Value(keyUserId)
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

	ctxUserID := ctx.Value(keyUserId)
	userId, ok := ctxUserID.(user.UserID)
	if !ok {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	rawWithdrawals, err :=s.balance.WithdrawalsHistory(ctx, userId)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if len(rawWithdrawals) == 0 {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	withdrawals := toWithdrawals(rawWithdrawals)

	result, err := json.Marshal(withdrawals)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func toWithdrawals(transactions []*balance.Transaction) []*withdrawal {
	countTransactions := len(transactions)
	res := make([]*withdrawal, countTransactions)

	for i := 0; i < countTransactions; i++ {
		res[i] = toWithdrawal(transactions[i])
	}

	return res
}

func toWithdrawal(transaction *balance.Transaction) *withdrawal {
	return &withdrawal{
		Order: strconv.FormatInt(transaction.OrderNum, 10),
		Sum: -transaction.Value,
		Processed_at: transaction.Processed_at,
	}
}
