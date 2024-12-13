package internal

import "time"

type Transaction struct {
	Id     int64
	User   string
	Amount float64
	Date   time.Time
}
