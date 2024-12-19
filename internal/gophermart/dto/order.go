package dto

import (
	"musthave_tpl/internal"
	"strconv"
)

type Order struct {
	Id         string  `json:"number"`
	Status     string  `json:"status"`
	Accrual    float64 `json:"accrual,omitempty"`
	UploadedAt string  `json:"uploaded_at"`
}

func NewOrders(orders []internal.Order) []Order {
	dto_orders := make([]Order, len(orders))
	for i, order := range orders {
		dto_orders[i] = Order{
			Id:         strconv.FormatInt(order.Number, 10),
			Status:     string(order.Status),
			Accrual:    order.Accrual,
			UploadedAt: order.UploadedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}
	return dto_orders
}
