package models

type OrderStatus string

const (
	New        OrderStatus = "NEW"
	Processing OrderStatus = "PROCESSING"
	Invalid    OrderStatus = "INVALID"
	Processed  OrderStatus = "PROCESSED"
)

type Order struct {
	Id      int64
	UserId  UserID
	Accrual float64
	Status  OrderStatus
}
