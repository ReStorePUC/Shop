package entity

import "time"

// Payment represents data about an payment.
type Payment struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Total     float64   `json:"total"`
	PIX       string    `json:"pix"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	StoreID   int       `json:"store_id"`
	ProductID int       `json:"product_id"`
}
