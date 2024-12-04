package internal

import "time"

type Transaction struct {
	Id     int64
	Amount int64
	Date   time.Time
}
