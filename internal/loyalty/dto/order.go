package dto

import (
	"musthave_tpl/internal/models"
	"strconv"
)

type Order struct {
	Order   string      `json:"order"`
	Status  OrderStatus `json:"status"`
	Accrual float64     `json:"accrual"`
}

type OrderStatus string

func (status OrderStatus) ToModel() models.OrderStatus {
	if status == Registered {
		return models.New
	}
	return models.OrderStatus(status)
}

const (
	Registered = "REGISTERED"
	Invalid    = "INVALID"
	Processing = "PROCESSING"
	Processed  = "PROCESSED"
)

func (o Order) ToInternal() (models.Order, error) {
	id, err := strconv.ParseInt(o.Order, 10, 64)
	return models.Order{
		Number:  id,
		Status:  o.Status.ToModel(),
		Accrual: o.Accrual,
	}, err
}
