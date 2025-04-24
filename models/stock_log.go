package models

import "time"

type StockLog struct {
	ID         int       `json:"id"`
	ProductID  int       `json:"product_id"`
	ChangeType string    `json:"change_type"` // "in" atau "out"
	Amount     int       `json:"amount"`
	Note       string    `json:"note"`
	CreatedAt  time.Time `json:"created_at"`
}
