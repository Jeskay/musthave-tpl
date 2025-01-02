package dto

import (
	"musthave_tpl/internal"
	"strconv"
)

type Order struct {
	Order   string      `json:"order"`
	Status  OrderStatus `json:"status"`
	Accrual float64     `json:"accrual"`
}

type OrderStatus string

func (status OrderStatus) ToInternal() internal.OrderStatus {
	if status == Registered {
		return internal.New
	}
	return internal.OrderStatus(status)
}

const (
	Registered = "REGISTERED"
	Invalid    = "INVALID"
	Processing = "PROCESSING"
	Processed  = "PROCESSED"
)

func (o Order) ToInternal() (internal.Order, error) {
	id, err := strconv.ParseInt(o.Order, 10, 64)
	return internal.Order{
		Number:  id,
		Status:  o.Status.ToInternal(),
		Accrual: o.Accrual,
	}, err
}
