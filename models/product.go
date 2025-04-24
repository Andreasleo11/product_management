package models

import "time"

type Product struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	PurchasePrice float64   `json:"purchase_price"`
	SellPrice     float64   `json:"sell_price"`
	CategoryID    int       `json:"category_id"`
	Stock         int       `json:"stock"`
	MinStockAlert int       `json:"min_stock_alert"`
	ImageURL      string    `json:"image_url"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
