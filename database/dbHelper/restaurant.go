package dbHelper

import (
	"database/sql"
	"errors"
	"rms/database"
	"rms/models"
	"time"
)

func CreateRestaurant(name, email, createdBy, address, state, city, pinCode string, lat, lng float64) (string, error) {
	arguments := []interface{}{
		name,
		email,
		createdBy,
		address,
		state,
		city,
		pinCode,
		lat,
		lng,
	}
	// language=SQL
	SQL := `INSERT INTO restaurants(name, email, created_by, address, state, city, pin_code, lat, lng) VALUES ($1, TRIM(LOWER($2)), $3, $4, $5, $6, $7, $8, $9) RETURNING id`
	var restaurantID string
	if err := database.RMS.QueryRowx(SQL, arguments...).Scan(&restaurantID); err != nil {
		return "", err
	}
	return restaurantID, nil
}

func CreateDish(restaurantID, createdBy, name, description string, quantity, price, discount int64) (string, error) {
	arguments := []interface{}{
		restaurantID,
		quantity,
		price,
		discount,
		createdBy,
		name,
		description,
	}
	// language=SQL
	SQL := `INSERT INTO dishes(restaurants_id, quantity, price, discount, created_by, name, description) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	var dishID string
	if err := database.RMS.QueryRowx(SQL, arguments...).Scan(&dishID); err != nil {
		return "", err
	}
	return dishID, nil
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
		if errors.Is(err, sql.ErrNoRows) {
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
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func UpdateRestaurant(restaurantID, name, email, Address, State, City, PinCode string, Lat, Lng float64) error {
	arguments := []interface{}{
		name,
		email,
		Address,
		State,
		City,
		PinCode,
		Lat,
		Lng,
		restaurantID,
	}
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
	_, err := database.RMS.Exec(SQL, arguments...)
	return err
}

func CloseMyRestaurant(restaurantID, createdBy string) error {
	//todo := dont return id **DONE**
	// language=SQL
	SQL := `UPDATE restaurants 
		SET archived_at = $1
		WHERE id = $2 AND created_by = $3`
	_, err := database.RMS.Exec(SQL, time.Now(), restaurantID, createdBy)
	return err
}

func CloseRestaurant(restaurantID string) error {
	// language=SQL
	SQL := `UPDATE restaurants 
		SET archived_at = $1
		WHERE id = $2`
	_, err := database.RMS.Exec(SQL, time.Now(), restaurantID)
	return err
}

func UpdateDish(dishID, restaurantId, name, description string, quantity, price, discount int64) error {
	arguments := []interface{}{
		name,
		description,
		quantity,
		price,
		discount,
		dishID,
		restaurantId,
	}
	// language=SQL
	SQL := `UPDATE dishes
		SET name = $1,
			description = $2,
			quantity = $3, 
			price = $4,
			discount = $5
		WHERE id = $6 AND restaurants_id = $7`
	_, err := database.RMS.Exec(SQL, arguments...)
	return err
}

func RemoveDishByUserID(dishID, restaurantID, createdBy string) error {
	// language=SQL
	SQL := `UPDATE dishes 
		SET archived_at = $1
		WHERE id = $2 AND restaurants_id = $3 AND created_by = $4`
	_, err := database.RMS.Exec(SQL, time.Now(), dishID, restaurantID, createdBy)
	return err
}

func RemoveDish(dishID, restaurantID string) error {
	// language=SQL
	SQL := `UPDATE dishes 
		SET archived_at = $1
		WHERE id = $2 AND restaurants_id = $3`
	_, err := database.RMS.Exec(SQL, time.Now(), dishID, restaurantID)
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &Dish, nil
}

func GetRestaurantsCountByUserID(createdBy string, Filters models.Filters) (int64, error) {
	// language=SQL
	SQL := `SELECT 
       			COUNT(r.id)
			FROM restaurants r
			WHERE r.archived_at IS NULL AND r.created_by = $1 AND
				r.name ILIKE '%' || $2 || '%' AND  r.email ILIKE '%' || $3 || '%'`
	var count int64
	err := database.RMS.Get(&count, SQL, createdBy, Filters.Name, Filters.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	return count, nil
}

func GetRestaurantsByUserID(createdBy string, Filters models.Filters) ([]models.Restaurant, error) {
	arguments := []interface{}{
		createdBy,
		Filters.Name,
		Filters.Email,
		Filters.SortBy,
		Filters.PageSize,
		Filters.PageSize * Filters.PageNumber,
	}
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
			WHERE r.archived_at IS NULL AND r.created_by = $1 AND
				r.name ILIKE '%' || $2 || '%' AND  r.email ILIKE '%' || $3 || '%'
			ORDER BY $4
			LIMIT $5
			OFFSET $6`
	restaurant := make([]models.Restaurant, 0)
	err := database.RMS.Select(&restaurant, SQL, arguments...)
	if err != nil {
		return nil, err
	}
	return restaurant, nil
}

func GetRestaurantByID(restaurantId string) (*models.Restaurant, error) {
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
	var restaurant models.Restaurant
	err := database.RMS.Get(&restaurant, SQL, restaurantId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &restaurant, nil
}

func GetRestaurantByIDAndUserID(restaurantId, createdBy string) (*models.Restaurant, error) {
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
			WHERE r.archived_at IS NULL AND r.restaurants_id = $1 AND r.created_by = $2`
	var restaurant models.Restaurant
	err := database.RMS.Get(&restaurant, SQL, restaurantId, createdBy)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &restaurant, nil
}

func GetRestaurantDishesCount(restaurantID string, Filters models.DishFilters) (int64, error) {
	arguments := []interface{}{
		restaurantID,
		Filters.CreatedBy,
		Filters.Name,
		Filters.MinQuantity,
		Filters.MinPrice,
		Filters.MaxPrice,
		Filters.MinDiscount,
		Filters.MaxDiscount,
	}
	// language=SQL
	SQL := `SELECT 
       			COUNT(d.id)
			FROM dishes d
			WHERE d.archived_at IS NULL AND d.restaurants_id = $1 AND d.created_by::text ILIKE '%' || $2 || '%' AND
			d.name ILIKE '%' || $3 || '%' AND  d.quantity > $4 AND d.price BETWEEN $5 AND $6 AND d.discount BETWEEN $7 AND $8`
	var count int64
	err := database.RMS.Get(&count, SQL, arguments...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	return count, nil
}

func GetRestaurantDishes(restaurantID string, Filters models.DishFilters) ([]models.Dishes, error) {
	arguments := []interface{}{
		restaurantID,
		Filters.CreatedBy,
		Filters.Name,
		Filters.MinQuantity,
		Filters.MinPrice,
		Filters.MaxPrice,
		Filters.MinDiscount,
		Filters.MaxDiscount,
		Filters.SortBy,
		Filters.PageSize,
		Filters.PageSize * Filters.PageNumber,
	}
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
			WHERE d.archived_at IS NULL AND d.restaurants_id = $1 AND d.created_by::text ILIKE '%' || $2 || '%' AND
			d.name ILIKE '%' || $3 || '%' AND  d.quantity > $4 AND d.price BETWEEN $5 AND $6 AND d.discount BETWEEN $7 AND $8
			ORDER BY $9
			LIMIT $10
			OFFSET $11`
	dishes := make([]models.Dishes, 0)
	err := database.RMS.Select(&dishes, SQL, arguments...)
	if err != nil {
		return nil, err
	}
	return dishes, nil
}

func GetRestaurantDishById(restaurantID, dishID string) (*models.Dishes, error) {
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
			WHERE d.archived_at IS NULL AND d.restaurants_id = $1 AND d.id = $2`
	var dishes models.Dishes
	err := database.RMS.Select(&dishes, SQL, restaurantID, dishID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &dishes, nil
}

func GetRestaurantDishesCountByUserId(restaurantID, createdBy string, Filters models.DishFilters) (int64, error) {
	arguments := []interface{}{
		restaurantID,
		createdBy,
		Filters.Name,
		Filters.MinQuantity,
		Filters.MinPrice,
		Filters.MaxPrice,
		Filters.MinDiscount,
		Filters.MaxDiscount,
	}
	// language=SQL
	SQL := `SELECT 
       			COUNT(d.id)
			FROM dishes d
			WHERE d.archived_at IS NULL AND d.restaurants_id = $1 AND d.created_by = $2 AND 
			d.name ILIKE '%' || $3 || '%' AND  d.quantity > $4 AND d.price BETWEEN $5 AND $6 AND d.discount BETWEEN $7 AND $8`
	var count int64
	err := database.RMS.Get(&count, SQL, arguments...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	return count, nil
}

func GetRestaurantDishesByUserID(restaurantID, createdBy string, Filters models.DishFilters) ([]models.Dishes, error) {
	arguments := []interface{}{
		restaurantID,
		createdBy,
		Filters.Name,
		Filters.MinQuantity,
		Filters.MinPrice,
		Filters.MaxPrice,
		Filters.MinDiscount,
		Filters.MaxDiscount,
		Filters.SortBy,
		Filters.PageSize,
		Filters.PageSize * Filters.PageNumber,
	}
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
			WHERE d.archived_at IS NULL AND d.restaurants_id = $1 AND d.created_by = $2 AND 
			d.name ILIKE '%' || $3 || '%' AND  d.quantity > $4 AND d.price BETWEEN $5 AND $6 AND d.discount BETWEEN $7 AND $8
			ORDER BY $9
			LIMIT $10
			OFFSET $11`
	dishes := make([]models.Dishes, 0)
	err := database.RMS.Select(&dishes, SQL, arguments...)
	if err != nil {
		return nil, err
	}
	return dishes, nil
}

func GetRestaurantDishByIDAndUserID(restaurantID, dishID, createdBy string) (*models.Dishes, error) {
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
			WHERE d.archived_at IS NULL AND d.restaurants_id = $1 AND d.id = $2 AND d.created_by = $3`
	var dishes models.Dishes
	err := database.RMS.Select(&dishes, SQL, restaurantID, dishID, createdBy)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &dishes, nil
}

func GetRestaurantsCount(Filters models.Filters) (int64, error) {
	// language=SQL
	//todo :=  use ilike because it searches matches like case insensetive matching **DONE**
	SQL := `SELECT 
       			COUNT(r.id)
			FROM restaurants r
			WHERE r.archived_at IS NULL AND r.created_by::text ILIKE '%' || $1 || '%'  AND
				r.name ILIKE '%' || $2 || '%' AND  r.email ILIKE '%' || $3 || '%'`
	var count int64
	err := database.RMS.Get(&count, SQL, Filters.CreatedBy, Filters.Name, Filters.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	return count, nil
}

func GetRestaurants(Filters models.Filters) ([]models.Restaurant, error) {
	arguments := []interface{}{
		Filters.CreatedBy,
		Filters.Name,
		Filters.Email,
		Filters.SortBy,
		Filters.PageSize,
		Filters.PageNumber * Filters.PageSize,
	}
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
			WHERE r.archived_at IS NULL AND r.created_by::text ILIKE '%' || $1 || '%'  AND
				r.name ILIKE '%' || $2 || '%' AND  r.email ILIKE '%' || $3 || '%'
			ORDER BY $4
			LIMIT $5
			OFFSET $6`
	restaurants := make([]models.Restaurant, 0)
	err := database.RMS.Select(&restaurants, SQL, arguments...)
	if err != nil {
		return nil, err
	}
	return restaurants, nil
}
