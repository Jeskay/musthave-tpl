package models

import "time"

type Transaction struct {
	ID     int64
	User   string
	Amount float64
	Date   time.Time
}
