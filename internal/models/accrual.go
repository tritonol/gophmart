package models

type Accrual struct {
	OrderNum string  `json:"order"`
	Status   string  `json:"status"`
	Value    float64 `json:"accrual"`
}
