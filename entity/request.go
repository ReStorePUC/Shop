package entity

import "time"

// Request represents data about an request.
type Request struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	PaymentID string    `json:"payment_id"`
	Price     float64   `json:"price"`
	Tax       float64   `json:"tax"`
	Track     string    `json:"track"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	StoreID   int       `json:"store_id"`
	ProductID int       `json:"product_id"`
}
