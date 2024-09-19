package handler

import (
	"net/http"
	"rms/database"
	"rms/database/dbHelper"
	"rms/middlewares"
	"rms/models"
	"rms/utils"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	body := struct {
		Name          string               `json:"name"`
		Email         string               `json:"email"`
		Password      string               `json:"password"`
		UserAddresses []models.UserAddress `json:"addresses"`
	}{}
	subAdminCtx := middlewares.UserContext(r)
	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.RespondError(w, http.StatusBadRequest, parseErr, "failed to parse request body")
		return
	}
	if len(body.Password) < 6 {
		utils.RespondError(w, http.StatusBadRequest, nil, "password must be 6 chars long")
		return
	}

	if !utils.IsEmailValid(body.Email) {
		utils.RespondError(w, http.StatusBadRequest, nil, "Invalid Email.")
		return
	}

	exists, existsErr := dbHelper.IsUserRoleExists(body.Email, models.RoleUser)
	if existsErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, existsErr, "failed to check user role existence")
		return
	}
	if exists {
		utils.RespondError(w, http.StatusBadRequest, nil, "user already exists")
		return
	}
	hashedPassword, hasErr := utils.HashPassword(body.Password)
	if hasErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, hasErr, "failed to secure password")
		return
	}
	txErr := database.Tx(func(tx *sqlx.Tx) error {
		userID, existsErr := dbHelper.IsUserExists(body.Email)
		if existsErr != nil {
			utils.RespondError(w, http.StatusInternalServerError, existsErr, "failed to check user existence")
			return existsErr
		}
		if len(userID) > 0 {
			roleErr := dbHelper.CreateUserRole(tx, userID, subAdminCtx.ID, models.RoleUser)
			if roleErr != nil {
				return roleErr
			}
		} else {
			userID, saveErr := dbHelper.CreateUser(tx, body.Name, body.Email, hashedPassword)
			if saveErr != nil {
				return saveErr
			}
			roleErr := dbHelper.CreateUserRole(tx, userID, subAdminCtx.ID, models.RoleUser)
			if roleErr != nil {
				return roleErr
			}
		}
		return nil
	})
	if txErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, txErr, "failed to create user")
		return
	}
	utils.RespondJSON(w, http.StatusCreated, struct {
		Message string `json:"message"`
	}{
		Message: "user Created successfully",
	})
}

func AddUserAddress(w http.ResponseWriter, r *http.Request) {
	body := struct {
		UserID  string  `json:"userId"`
		Address string  `json:"address"`
		State   string  `json:"state"`
		City    string  `json:"city"`
		PinCode string  `json:"pinCode"`
		Lat     float64 `json:"lat"`
		Lng     float64 `json:"lng"`
	}{}
	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.RespondError(w, http.StatusBadRequest, parseErr, "failed to parse request body")
		return
	}

	exists, existsErr := dbHelper.IsUserRoleWithUserIDExists(body.UserID, models.RoleUser)
	if existsErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, existsErr, "failed to check user role existence")
		return
	}
	if !exists {
		utils.RespondError(w, http.StatusBadRequest, nil, "user not exists")
		return
	}

	if len(body.Address) > 30 || len(body.Address) <= 2 {
		utils.RespondError(w, http.StatusBadRequest, nil, "Address must be with in 2 to 30 letter.")
		return
	}

	if len(body.State) > 16 || len(body.State) <= 2 {
		utils.RespondError(w, http.StatusBadRequest, nil, "State must be with in 2 to 16 letter.")
		return
	}

	if len(body.City) > 20 || len(body.City) <= 2 {
		utils.RespondError(w, http.StatusBadRequest, nil, "City must be with in 2 to 20 letter.")
		return
	}

	if len(body.City) > 20 || len(body.City) <= 2 {
		utils.RespondError(w, http.StatusBadRequest, nil, "City must be with in 2 to 20 letter.")
		return
	}

	if body.Lat > 90 || body.Lat < -90 {
		utils.RespondError(w, http.StatusBadRequest, nil, "PinCode must 6 digit.")
		return
	}

	if body.Lng > 180 || body.Lng < -180 {
		utils.RespondError(w, http.StatusBadRequest, nil, "PinCode must 6 digit.")
		return
	}

	txErr := database.Tx(func(tx *sqlx.Tx) error {
		addressErr := dbHelper.CreateUserAddress(tx, body.UserID, body.Address, body.State, body.City, body.PinCode, body.Lat, body.Lng)
		if addressErr != nil {
			utils.RespondError(w, http.StatusBadRequest, addressErr, "Invalid Address.")
			return addressErr
		}
		return nil
	})
	if txErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, txErr, "failed to create user")
		return
	}
	utils.RespondJSON(w, http.StatusCreated, struct {
		Message string `json:"message"`
	}{
		Message: "Address Added to user successfully",
	})
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	adminCtx := middlewares.UserContext(r)
	Users, err := dbHelper.GetUsers(adminCtx.ID, models.RoleUser)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "Unable to get Users")
		return
	}
	utils.RespondJSON(w, http.StatusCreated, struct {
		Message string        `json:"message"`
		Users   []models.User `json:"users"`
	}{
		Message: "Get users successfully.",
		Users:   Users,
	})
}

func RemoveMyUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "userId")
	adminCtx := middlewares.UserContext(r)
	err := dbHelper.RemoveMyUser(id, adminCtx.ID, models.RoleUser)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "Unable to get Users")
		return
	}
	utils.RespondJSON(w, http.StatusCreated, struct {
		Message string `json:"message"`
	}{
		Message: "User remove successfully.",
	})
}

// Restaurant

func OpenRestaurant(w http.ResponseWriter, r *http.Request) {
	body := struct {
		Name    string  `json:"name"`
		Email   string  `json:"email"`
		Address string  `json:"address"`
		State   string  `json:"state"`
		City    string  `json:"city"`
		PinCode string  `json:"pinCode"`
		Lat     float64 `json:"lat"`
		Lng     float64 `json:"lng"`
	}{}
	adminCtx := middlewares.UserContext(r)
	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.RespondError(w, http.StatusBadRequest, parseErr, "failed to parse request body")
		return
	}

	if !utils.IsEmailValid(body.Email) {
		utils.RespondError(w, http.StatusBadRequest, nil, "Invalid Email.")
		return
	}

	exists, existsErr := dbHelper.IsRestaurantExists(body.Email)
	if existsErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, existsErr, "failed to check Restaurant existence")
		return
	}

	if exists {
		utils.RespondError(w, http.StatusBadRequest, nil, "Restaurant already exists")
		return
	}

	if body.Name == "" {
		utils.RespondError(w, http.StatusBadRequest, nil, "Invalid Name")
		return
	}
	if !utils.IsEmailValid(body.Email) {
		utils.RespondError(w, http.StatusBadRequest, nil, "Invalid Email")
		return
	}

	if len(body.Address) > 30 || len(body.Address) <= 2 {
		utils.RespondError(w, http.StatusBadRequest, nil, "Address must be with in 2 to 30 letter.")
		return
	}

	if len(body.State) > 16 || len(body.State) <= 2 {
		utils.RespondError(w, http.StatusBadRequest, nil, "State must be with in 2 to 16 letter.")
		return
	}

	if len(body.City) > 20 || len(body.City) <= 2 {
		utils.RespondError(w, http.StatusBadRequest, nil, "City must be with in 2 to 20 letter.")
		return
	}

	if len(body.PinCode) != 6 {
		utils.RespondError(w, http.StatusBadRequest, nil, "PinCode must 6 digit.")
		return
	}

	if body.Lat > 90 || body.Lat < -90 {
		utils.RespondError(w, http.StatusBadRequest, nil, "Invalid Latitude.")
		return
	}

	if body.Lng > 180 || body.Lng < -180 {
		utils.RespondError(w, http.StatusBadRequest, nil, "Invalid Longitude.")
		return
	}

	txErr := database.Tx(func(tx *sqlx.Tx) error {
		_, saveErr := dbHelper.CreateRestaurant(tx, body.Name, body.Email, adminCtx.ID, body.Address, body.State, body.City, body.PinCode, body.Lat, body.Lng)
		if saveErr != nil {
			return saveErr
		}
		return nil
	})

	if txErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, txErr, "failed to create user")
		return
	}
	utils.RespondJSON(w, http.StatusCreated, struct {
		Message string `json:"message"`
	}{
		Message: "Restaurant Opened successfully",
	})
}

func CloseMyRestaurant(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "restaurantId")
	adminCtx := middlewares.UserContext(r)
	err := dbHelper.CloseMyRestaurant(id, adminCtx.ID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "Unable to get Users")
		return
	}
	utils.RespondJSON(w, http.StatusCreated, struct {
		Message string `json:"message"`
	}{
		Message: "Restaurant Closed successfully.",
	})
}

func GetMyRestaurants(w http.ResponseWriter, r *http.Request) {
	adminCtx := middlewares.UserContext(r)
	Restaurants, err := dbHelper.GetRestaurantsByUserID(adminCtx.ID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "Unable to get Users")
		return
	}
	utils.RespondJSON(w, http.StatusCreated, struct {
		Message     string              `json:"message"`
		Restaurants []models.Restaurant `json:"restaurants"`
	}{
		Message:     "Get Restaurants successfully.",
		Restaurants: Restaurants,
	})
}

func AddRestaurantDish(w http.ResponseWriter, r *http.Request) {
	body := struct {
		Name         string `json:"name"`
		Description  string `json:"description"`
		RestaurantID string `json:"restaurantID"`
		Quantity     int64  `json:"quantity"`
		Price        int64  `json:"price"`
		Discount     int64  `json:"discount"`
	}{}
	adminCtx := middlewares.UserContext(r)
	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.RespondError(w, http.StatusBadRequest, parseErr, "failed to parse request body")
		return
	}

	exists, existsErr := dbHelper.IsRestaurantIDExists(body.RestaurantID)
	if existsErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, existsErr, "failed to check Restaurant existence")
		return
	}

	if !exists {
		utils.RespondError(w, http.StatusBadRequest, nil, "Restaurant not exists")
		return
	}

	if body.Name == "" {
		utils.RespondError(w, http.StatusBadRequest, nil, "Invalid Dish Name.")
		return
	}

	if body.Description == "" {
		utils.RespondError(w, http.StatusBadRequest, nil, "Invalid Dish Description.")
		return
	}

	if body.Quantity <= 0 {
		utils.RespondError(w, http.StatusBadRequest, nil, "Invalid Dish Quantity.")
		return
	}

	if body.Price <= 0 {
		utils.RespondError(w, http.StatusBadRequest, nil, "Invalid Dish Price.")
		return
	}

	if body.Discount > 100 || body.Discount < 0 {
		utils.RespondError(w, http.StatusBadRequest, nil, "Invalid Discount.")
		return
	}

	txErr := database.Tx(func(tx *sqlx.Tx) error {
		_, saveErr := dbHelper.CreateDish(tx, body.RestaurantID, adminCtx.ID, body.Name, body.Description, body.Quantity, body.Price, body.Discount)
		if saveErr != nil {
			return saveErr
		}
		return nil
	})

	if txErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, txErr, "failed to create Dish")
		return
	}
	utils.RespondJSON(w, http.StatusCreated, struct {
		Message string `json:"message"`
	}{
		Message: "Restaurant Dishes added successfully",
	})
}

func GetMyDishes(w http.ResponseWriter, r *http.Request) {
	adminCtx := middlewares.UserContext(r)
	Menu, err := dbHelper.GetDishesByUserID(adminCtx.ID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "Unable to get Users")
		return
	}
	utils.RespondJSON(w, http.StatusCreated, struct {
		Message string          `json:"message"`
		Menu    []models.Dishes `json:"menu"`
	}{
		Message: "Get Dishes successfully.",
		Menu:    Menu,
	})
}

func UpdateMyRestaurant(w http.ResponseWriter, r *http.Request) {
	body := struct {
		ID      string  `json:"id"`
		Name    string  `json:"name"`
		Email   string  `json:"email"`
		Address string  `json:"address"`
		State   string  `json:"state"`
		City    string  `json:"city"`
		PinCode string  `json:"pinCode"`
		Lat     float64 `json:"lat"`
		Lng     float64 `json:"lng"`
	}{}

	userCtx := middlewares.UserContext(r)
	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.RespondError(w, http.StatusBadRequest, parseErr, "failed to parse request body")
		return
	}

	restaurant, restaurantErr := dbHelper.GetRestaurantByID(body.ID)
	if restaurantErr != nil {
		utils.RespondError(w, http.StatusBadRequest, nil, "Restaurant not exist")
		return
	}

	if restaurant.CreatedBy != userCtx.ID {
		utils.RespondError(w, http.StatusBadRequest, nil, "Restaurant not Created by: "+string(userCtx.CurrentRole))
		return
	}

	if body.Name == "" {
		body.Name = restaurant.Name
	}
	if !utils.IsEmailValid(body.Email) {
		if body.Email != "" {
			utils.RespondError(w, http.StatusBadRequest, nil, "Invalid Email")
			return
		}
		body.Email = restaurant.Email
	}

	if len(body.Address) > 30 || len(body.Address) <= 2 {
		if body.Address != "" {
			utils.RespondError(w, http.StatusBadRequest, nil, "Address must be with in 2 to 30 letter.")
			return
		}
		body.Address = restaurant.Address
	}

	if len(body.State) > 16 || len(body.State) <= 2 {
		if body.State != "" {
			utils.RespondError(w, http.StatusBadRequest, nil, "State must be with in 2 to 16 letter.")
			return
		}
		body.State = restaurant.State
	}

	if len(body.City) > 20 || len(body.City) <= 2 {
		if body.City != "" {
			utils.RespondError(w, http.StatusBadRequest, nil, "City must be with in 2 to 20 letter.")
			return
		}
		body.City = restaurant.City
	}

	if len(body.PinCode) != 6 {
		if body.PinCode != "" {
			utils.RespondError(w, http.StatusBadRequest, nil, "PinCode must 6 digit.")
			return
		}
		body.PinCode = restaurant.PinCode
	}

	if body.Lat > 90 || body.Lat < -90 {
		if body.Lat != 0 {
			utils.RespondError(w, http.StatusBadRequest, nil, "Invalid Latitude.")
			return
		}
		body.Lat = restaurant.Lat
	}

	if body.Lng > 180 || body.Lng < -180 {
		if body.Lat != 0 {
			utils.RespondError(w, http.StatusBadRequest, nil, "Invalid Longitude.")
			return
		}
		body.Lng = restaurant.Lng
	}

	err := dbHelper.UpdateRestaurant(body.ID, body.Name, body.Email, body.Address, body.State, body.City, body.PinCode, body.Lat, body.Lng)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed update User")
		return
	}
	utils.RespondJSON(w, http.StatusCreated, struct {
		Message string `json:"message"`
	}{
		Message: "Address Update successfully",
	})
}

func UpdateMyDish(w http.ResponseWriter, r *http.Request) {
	body := struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Quantity    int64  `json:"quantity"`
		Price       int64  `json:"price"`
		Discount    int64  `json:"discount"`
	}{}

	userCtx := middlewares.UserContext(r)
	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.RespondError(w, http.StatusBadRequest, parseErr, "failed to parse request body")
		return
	}

	dish, dishErr := dbHelper.GetDishByID(body.ID)
	if dishErr != nil {
		utils.RespondError(w, http.StatusBadRequest, nil, "Dish not exist")
		return
	}

	if dish.CreatedBy != userCtx.ID {
		utils.RespondError(w, http.StatusBadRequest, nil, "Dish not Created by: "+string(userCtx.CurrentRole))
		return
	}

	if body.Name == "" {
		body.Name = dish.Name
	}

	if body.Description == "" {
		body.Name = dish.Description
	}

	if body.Quantity <= 0 {
		body.Quantity = dish.Quantity
	}

	if body.Price <= 0 {
		body.Price = dish.Price
	}

	if body.Discount > 100 || body.Discount < 0 {
		utils.RespondError(w, http.StatusBadRequest, nil, "Invalid Discount.")
		return
	}

	err := dbHelper.UpdateDish(body.ID, body.Name, body.Description, body.Quantity, body.Price, body.Discount)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed update User")
		return
	}
	utils.RespondJSON(w, http.StatusCreated, struct {
		Message string `json:"message"`
	}{
		Message: "Address Update successfully",
	})
}

func RemoveMyDish(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "dishId")
	adminCtx := middlewares.UserContext(r)
	err := dbHelper.RemoveMyDish(id, adminCtx.ID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "Unable to get Users")
		return
	}
	utils.RespondJSON(w, http.StatusCreated, struct {
		Message string `json:"message"`
	}{
		Message: "Dish Removed successfully.",
	})
}
