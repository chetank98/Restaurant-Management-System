package models

import "time"

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleSubAdmin Role = "sub-admin"
	RoleUser     Role = "user"
)

func (r Role) IsValid() bool {
	return r == RoleAdmin || r == RoleSubAdmin || r == RoleUser
}

type User struct {
	ID            string        `json:"id" db:"id"`
	Name          string        `json:"name" db:"name"`
	Email         string        `json:"email" db:"email"`
	Password      string        `json:"password" db:"password"`
	CreatedAt     time.Time     `json:"createdAt" db:"created_at"`
	CurrentRole   Role          `json:"currentRole" db:"user_current_role"`
	UserAddresses []UserAddress `json:"Addresses" db:"user_addresses"`
}

type UserAddress struct {
	ID        string    `json:"id" db:"id"`
	Address   string    `json:"address" db:"address"`
	State     string    `json:"state" db:"state"`
	City      string    `json:"city" db:"city"`
	PinCode   string    `json:"pinCode" db:"pin_code"`
	Lat       float64   `json:"lat" db:"lat"`
	Lng       float64   `json:"lng" db:"lng"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}
