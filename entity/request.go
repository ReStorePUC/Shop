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
	UserID    int       `json:"user_id"`
	Product   *Product  `json:"product"`
}

type Create struct {
	Items []Request `json:"items"`
}

// Product represents data about an product.
type Product struct {
	ID          int     `json:"id" gorm:"primaryKey"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Categories  string  `json:"categories"`
	Size        string  `json:"size"`
	Price       float64 `json:"price"`
	Tax         float64 `json:"tax"`
	Available   bool    `json:"available"`
	StoreID     int     `json:"store_id"`
	Images      []Image `json:"images"`
}

// Image represents data about an image.
type Image struct {
	ID        int    `json:"id" gorm:"primaryKey"`
	ImagePath string `json:"image_path"`
	ProductID int    `json:"product_id"`
}
