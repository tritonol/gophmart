package models

import "github.com/tritonol/gophmart.git/internal/models/user"

type OrderStatus string

const (
	New        OrderStatus = "NEW"
	Processing OrderStatus = "PROCESSING"
	Invalid    OrderStatus = "INVALID"
	Processed  OrderStatus = "PROCESSED"
)

type Order struct {
	Id      int64
	UserId  user.UserID
	Accrual float64
	Status  OrderStatus
}
