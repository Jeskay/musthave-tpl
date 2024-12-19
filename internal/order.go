package internal

import "time"

type Order struct {
	User       User
	Number     int64
	Status     OrderStatus
	Accrual    float64
	UploadedAt time.Time
}

type OrderStatus string

const (
	New        = "NEW"
	Invalid    = "INVALID"
	Processing = "PROCESSING"
	Processed  = "PROCESSED"
)
