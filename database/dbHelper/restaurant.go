package dbHelper

import (
	"database/sql"
	"rms/database"
	"rms/models"
	"time"

	"github.com/jmoiron/sqlx"
)

func CreateRestaurant(db sqlx.Ext, name, email, createdBy, address, state, city, pinCode string, lat, lng float64) (string, error) {
	// language=SQL
	SQL := `INSERT INTO restaurants(name, email, created_by, address, state, city, pin_code, lat, lng) VALUES ($1, TRIM(LOWER($2)), $3, $4, $5, $6, $7, $8, $9) RETURNING id`
	var userID string
	if err := db.QueryRowx(SQL, name, email, createdBy, address, state, city, pinCode, lat, lng).Scan(&userID); err != nil {
		return "", err
	}
	return userID, nil
}

func CreateDish(db sqlx.Ext, restaurantID, createdBy, name, description string, quantity, price, discount int64) (string, error) {
	// language=SQL
	SQL := `INSERT INTO dishes(restaurants_id, quantity, price, discount, created_by, name, description) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	var userID string
	if err := db.QueryRowx(SQL, restaurantID, quantity, price, discount, createdBy, name, description).Scan(&userID); err != nil {
		return "", err
	}
	return userID, nil
}

func IsRestaurantExists(email string) (bool, error) {
	// language=SQL
	SQL := `SELECT
				r.id
       		FROM
				restaurants r
			WHERE
				r.archived_at IS NULL
				AND r.email = TRIM(LOWER($1))`
	var restaurantId string
	err := database.RMS.Get(&restaurantId, SQL, email)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func IsRestaurantIDExists(id string) (bool, error) {
	// language=SQL
	SQL := `SELECT
				r.id
       		FROM
				restaurants r
			WHERE
				r.archived_at IS NULL
				AND r.id = $1`
	var restaurantId string
	err := database.RMS.Get(&restaurantId, SQL, id)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func GetRestaurantByID(restaurantID string) (*models.Restaurant, error) {
	// language=SQL
	SQL := `SELECT 
       			r.id,
       			r.name,
       			r.email,
       			r.created_at,
       			r.created_by,
				r.address,
				r.state,
				r.city,
				r.pin_code,
				r.lat,
				r.lng
			FROM restaurants r
			WHERE r.archived_at IS NULL AND r.id = $1`
	var Restaurant models.Restaurant
	err := database.RMS.Get(&Restaurant, SQL, restaurantID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &Restaurant, nil
}

func UpdateRestaurant(restaurantID, name, email, Address, State, City, PinCode string, Lat, Lng float64) error {
	// language=SQL
	SQL := `UPDATE restaurants
		SET name = $1,
			email = $2
			address = $3, 
			state = $4,
			city = $5, 
			pin_code = $6,
			lat = $7, 
			lng = $8
		WHERE id = $9
		RETURNING id;`
	_, err := database.RMS.Exec(SQL, name, email, Address, State, City, PinCode, Lat, Lng, restaurantID)
	return err
}

func CloseMyRestaurant(restaurantID, createdBy string) error {
	// language=SQL
	SQL := `UPDATE restaurants 
		SET archived_at = $1
		WHERE id = $2 AND created_by = $3
		RETURNING id;`
	_, err := database.RMS.Exec(SQL, time.Now(), restaurantID, createdBy)
	return err
}

func CloseRestaurant(restaurantID string) error {
	// language=SQL
	SQL := `UPDATE restaurants 
		SET archived_at = $1
		WHERE id = $2
		RETURNING id;`
	_, err := database.RMS.Exec(SQL, time.Now(), restaurantID)
	return err
}

func UpdateDish(dishID, name, description string, quantity, price, discount int64) error {
	// language=SQL
	SQL := `UPDATE dishes
		SET name = $1,
			description = $2
			quantity = $3, 
			price = $4,
			discount = $5
		WHERE id = $6
		RETURNING id;`
	_, err := database.RMS.Exec(SQL, name, description, quantity, price, discount, dishID)
	return err
}

func RemoveMyDish(dishID, createdBy string) error {
	// language=SQL
	SQL := `UPDATE dishes 
		SET archived_at = $1
		WHERE id = $2 AND created_by = $3
		RETURNING id;`
	_, err := database.RMS.Exec(SQL, time.Now(), dishID, createdBy)
	return err
}

func RemoveDish(dishID string) error {
	// language=SQL
	SQL := `UPDATE dishes 
		SET archived_at = $1
		WHERE id = $2
		RETURNING id;`
	_, err := database.RMS.Exec(SQL, time.Now(), dishID)
	return err
}

func GetDishByID(dishID string) (*models.Dishes, error) {
	// language=SQL
	SQL := `SELECT 
       			d.id,
       			d.name,
       			d.description,
				d.quantity,
				d.price,
				d.discount,
       			d.created_at,
       			d.created_by
			FROM dishes d
			WHERE d.archived_at IS NULL AND d.id = $1`
	var Dish models.Dishes
	err := database.RMS.Get(&Dish, SQL, dishID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &Dish, nil
}

func GetRestaurantsByUserID(createdBy string) ([]models.Restaurant, error) {
	// language=SQL
	SQL := `SELECT 
       			r.id,
       			r.name,
       			r.email,
       			r.created_at,
       			r.created_by,
				r.address,
				r.state,
				r.city,
				r.pin_code,
				r.lat,
				r.lng
			FROM restaurants r
			WHERE r.archived_at IS NULL AND r.created_by = $1`
	Restaurant := make([]models.Restaurant, 0)
	err := database.RMS.Select(&Restaurant, SQL, createdBy)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return Restaurant, nil
}

func GetRestaurantDish(restaurantID string) ([]models.Dishes, error) {
	// language=SQL
	SQL := `SELECT 
       			d.id,
       			d.name,
       			d.description,
				d.quantity,
				d.price,
				d.discount,
       			d.created_at,
       			d.created_by
			FROM dishes d
			WHERE d.archived_at IS NULL AND d.restaurants_id = $1`
	Dishes := make([]models.Dishes, 0)
	err := database.RMS.Select(&Dishes, SQL, restaurantID)
	if err != nil {
		return nil, err
	}
	return Dishes, nil
}

func GetDishesByUserID(createdBy string) ([]models.Dishes, error) {
	// language=SQL
	SQL := `SELECT 
       			d.id,
       			d.name,
       			d.description,
				d.quantity,
				d.price,
				d.discount,
       			d.created_at,
       			d.created_by
			FROM dishes d
			WHERE d.archived_at IS NULL AND d.created_by = $1`
	Dishes := make([]models.Dishes, 0)
	err := database.RMS.Select(&Dishes, SQL, createdBy)
	if err != nil {
		return nil, err
	}
	return Dishes, nil
}

func GetAllDishes() ([]models.Dishes, error) {
	// language=SQL
	SQL := `SELECT 
       			d.id,
       			d.name,
       			d.description,
				d.quantity,
				d.price,
				d.discount,
       			d.created_at,
       			d.created_by
			FROM dishes d
			WHERE d.archived_at IS NULL`
	Dishes := make([]models.Dishes, 0)
	err := database.RMS.Select(&Dishes, SQL)
	if err != nil {
		return nil, err
	}
	return Dishes, nil
}

func GetFilteredDishes(name, sortedBy string, Range int64) ([]models.Dishes, error) {
	// language=SQL
	SQL := `SELECT 
       			d.id,
       			d.name,
       			d.description,
				d.quantity,
				d.price,
				d.discount,
       			d.created_at,
       			d.created_by
			FROM dishes d
			WHERE d.archived_at IS NULL AND d.name LIKE '%' || $1 || '%'
			ORDER BY $2
			LIMIT $3`
	Dishes := make([]models.Dishes, 0)
	err := database.RMS.Select(&Dishes, SQL, name, sortedBy, Range)
	if err != nil {
		return nil, err
	}
	return Dishes, nil
}

func GetAllRestaurants() ([]models.Restaurant, error) {
	// language=SQL
	SQL := `SELECT 
       			r.id,
       			r.name,
       			r.email,
       			r.created_at,
       			r.created_by,
				r.address,
				r.state,
				r.city,
				r.pin_code,
				r.lat,
				r.lng
			FROM restaurants r
			WHERE r.archived_at IS NULL`
	Restaurant := make([]models.Restaurant, 0)
	err := database.RMS.Select(&Restaurant, SQL)
	if err != nil {
		return nil, err
	}
	return Restaurant, nil
}

func GetFilteredRestaurants(name, sortedBy string, Range int64) ([]models.Restaurant, error) {
	// language=SQL
	SQL := `SELECT 
       			r.id,
       			r.name,
       			r.email,
       			r.created_at,
       			r.created_by,
				r.address,
				r.state,
				r.city,
				r.pin_code,
				r.lat,
				r.lng
			FROM restaurants r
			WHERE r.archived_at IS NULL AND r.name LIKE '%' || $1 || '%'
			ORDER BY $2
			LIMIT $3`
	Restaurant := make([]models.Restaurant, 0)
	err := database.RMS.Select(&Restaurant, SQL, name, sortedBy, Range)
	if err != nil {
		return nil, err
	}
	return Restaurant, nil
}
