package models

import "time"

// Restaurant

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

type OpenRestaurantBody struct {
	Name      string  `json:"name" db:"name"`
	Email     string  `json:"email" db:"email"`
	Address   string  `json:"address" db:"address"`
	State     string  `json:"state" db:"state"`
	City      string  `json:"city" db:"city"`
	PinCode   string  `json:"pinCode" db:"pin_code"`
	Lat       float64 `json:"lat" db:"lat"`
	Lng       float64 `json:"lng" db:"lng"`
	CreatedBy string  `json:"createdBy" db:"created_by"`
}

type GetRestaurants struct {
	Message     string       `json:"message"`
	Restaurants []Restaurant `json:"restaurants"`
	TotalCount  int64        `json:"totalCount"`
	PageNumber  int64        `json:"pageNumber"`
	PageSize    int64        `json:"pageSize"`
}

type RestaurantDistance struct {
	Message      string  `json:"message"`
	Distance     float64 `json:"restaurantsDistance"`
	DistanceUnit string  `json:"distanceUnit"`
}

// Dishes

type DishSortedBy string

const (
	DishID        DishSortedBy = "id"
	DishName      DishSortedBy = "name"
	DishQuantity  DishSortedBy = "quantity"
	DishPrice     DishSortedBy = "price"
	DishDiscount  DishSortedBy = "discount"
	DishCreatedBy DishSortedBy = "created_by"
)

func (ds DishSortedBy) IsValid() bool {
	return ds == DishID || ds == DishName || ds == DishCreatedBy
}

type DishFilters struct {
	PageNumber  int64
	PageSize    int64
	Name        string
	Email       string
	MinQuantity int64
	MaxPrice    int64
	MinPrice    int64
	MaxDiscount int64
	MinDiscount int64
	CreatedBy   string
	SortBy      DishSortedBy
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

type AddDishesBody struct {
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	Quantity    int64  `json:"quantity" db:"quantity"`
	Price       int64  `json:"price" db:"price"`
	Discount    int64  `json:"discount" db:"discount"`
	CreatedBy   string `json:"createdBy" db:"created_by"`
}

type GetDishes struct {
	Message    string   `json:"message"`
	Dishes     []Dishes `json:"dishes"`
	TotalCount int64    `json:"totalCount"`
	PageNumber int64    `json:"pageNumber"`
	PageSize   int64    `json:"pageSize"`
}
