package models

import "time"

// Transaction - withdrawal of given amount of bonuses from user balance in favor of given order.
type Transaction struct {
	ID     int64
	User   string
	Amount float64
	Date   time.Time
}
