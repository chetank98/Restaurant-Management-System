package models

import "time"

type Restaurant struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	Address   string    `json:"address" db:"address"`
	State     string    `json:"state" db:"state"`
	City      string    `json:"city" db:"city"`
	PinCode   string    `json:"pinCode" db:"pin_code"`
	Lat       float64   `json:"lat" db:"lat"`
	Lng       float64   `json:"lng" db:"lng"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	CreatedBy string    `json:"createdBy" db:"created_by"`
}

type Dishes struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Quantity    int64     `json:"quantity" db:"quantity"`
	Price       int64     `json:"price" db:"price"`
	Discount    int64     `json:"discount" db:"discount"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	CreatedBy   string    `json:"createdBy" db:"created_by"`
}
