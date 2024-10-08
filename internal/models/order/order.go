package order

import "github.com/tritonol/gophmart.git/internal/models/user"

type OrderStatus string

const (
	New        OrderStatus = "NEW"
	Processing OrderStatus = "PROCESSING"
	Invalid    OrderStatus = "INVALID"
	Processed  OrderStatus = "PROCESSED"
)

type Order struct {
	ID         int64       `json:"number"`
	Accrual    float64     `json:"accrual"`
	Status     OrderStatus `json:"status"`
	UploadedAt string      `json:"uploaded_at"`
	UserID     user.UserID
}
