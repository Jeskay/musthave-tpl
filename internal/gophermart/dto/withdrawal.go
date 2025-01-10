package dto

import (
	"musthave_tpl/internal/models"
	"strconv"
)

type Withdrawal struct {
	OrderID     string  `json:"order"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at,omitempty"`
}

func NewWithdrawals(transactions []models.Transaction) []Withdrawal {
	withdrawals := make([]Withdrawal, len(transactions))
	for i, t := range transactions {
		withdrawals[i] = Withdrawal{
			OrderID:     strconv.FormatInt(t.ID, 10),
			Sum:         t.Amount,
			ProcessedAt: t.Date.Format("2006-01-02T15:04:05Z07:00"),
		}
	}
	return withdrawals
}
