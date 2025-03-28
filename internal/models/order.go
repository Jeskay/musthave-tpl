package models

import "time"

// Order - order, registered in gophermart loyalty system.
type Order struct {
	User       User
	Number     int64
	Status     OrderStatus
	Accrual    float64
	UploadedAt time.Time
}

// OrderStatus - current status of Order in gophermart system.
type OrderStatus string

const (
	// New - Order has been registered in he system and awaits for processing.
	New = "NEW"
	// Invalid - Loyalty system denied accrual for the order.
	Invalid = "INVALID"
	// Processing - accrual of bonus in process.
	Processing = "PROCESSING"
	// Processed - Order has been verified, accrual bonus information has been received.
	Processed = "PROCESSED"
)
