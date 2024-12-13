package dto

import (
	"musthave_tpl/internal"
	"strconv"
)

type Order struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

func (o Order) ToInternal() (internal.Order, error) {
	id, err := strconv.ParseInt(o.Order, 10, 64)
	return internal.Order{
		Number:  id,
		Status:  internal.OrderStatus(o.Status),
		Accrual: o.Accrual,
	}, err
}
