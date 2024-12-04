package internal

type Order struct {
	User    User
	Number  int64
	Status  OrderStatus
	Accrual float64
}

type OrderStatus string

const (
	Registered = "REGISTERED"
	Invalid    = "INVALID"
	Processing = "PROCESSING"
	Processed  = "PROCESSED"
)
