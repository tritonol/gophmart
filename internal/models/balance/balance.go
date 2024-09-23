package balance

import "github.com/tritonol/gophmart.git/internal/models/user"

type Transaction struct {
	ID           int64       `db:"id"`
	OrderNum     int64       `db:"from_id"`
	UserID       user.UserID `db:"user_id"`
	Value        float64     `db:"value"`
	ProcessedAt string      `db:"processed_at"`
}

type Balance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}