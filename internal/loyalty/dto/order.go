// Module dto describes objects used in client-server communication.
package dto

import (
	"musthave_tpl/internal/models"
	"strconv"
)

// Order - order, registered in loyalty server. Describes JSON representation.
type Order struct {
	Order   string      `json:"order"`
	Status  OrderStatus `json:"status"`
	Accrual float64     `json:"accrual"`
}

// OrderStatus - current processing status of Order in loyalty system.
type OrderStatus string

// ToModel converts JSON Status type to internal status type.
func (status OrderStatus) ToModel() models.OrderStatus {
	if status == Registered {
		return models.New
	}
	return models.OrderStatus(status)
}

const (
	// Registered - Order has been registered in the system but accrual is not yet calculated.
	Registered = "REGISTERED"
	// Invalid - Order accrual denied.
	Invalid = "INVALID"
	// Processing - accrual of bonus in process.
	Processing = "PROCESSING"
	// Processed - accruing process complete.
	Processed = "PROCESSED"
)

// ToInternal converts JSON Order type to internal Order type.
func (o Order) ToInternal() (models.Order, error) {
	id, err := strconv.ParseInt(o.Order, 10, 64)
	return models.Order{
		Number:  id,
		Status:  o.Status.ToModel(),
		Accrual: o.Accrual,
	}, err
}
