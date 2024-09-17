package balance

import "github.com/tritonol/gophmart.git/internal/models/user"

type Transaction struct {
	Id           int64       `db:"id"`
	OrderNum     int64       `db:"from_id"`
	UserId       user.UserID `db:"user_id"`
	Value        float64     `db:"value"`
	Processed_at string      `db:"processed_at"`
}
