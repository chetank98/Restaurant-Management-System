package handler

import (
	"net/http"
	"os"
	"rms/database"
	"rms/database/dbHelper"
	"rms/middlewares"
	"rms/models"
	"rms/utils"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

func RegisterSubAdmin(w http.ResponseWriter, r *http.Request) {
	body := struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	adminCtx := middlewares.UserContext(r)
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

	exists, existsErr := dbHelper.IsUserRoleExists(body.Email, models.RoleSubAdmin)
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
			roleErr := dbHelper.CreateUserRole(tx, userID, adminCtx.ID, models.RoleSubAdmin)
			if roleErr != nil {
				return roleErr
			}
		} else {
			userID, saveErr := dbHelper.CreateUser(tx, body.Name, body.Email, hashedPassword)
			if saveErr != nil {
				return saveErr
			}
			roleErr := dbHelper.CreateUserRole(tx, userID, adminCtx.ID, models.RoleSubAdmin)
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
		Message: "SubAdmin Created successfully",
	})
}

func GetSubAdmins(w http.ResponseWriter, r *http.Request) {
	adminCtx := middlewares.UserContext(r)
	subAdmins, err := dbHelper.GetUsers(adminCtx.ID, models.RoleSubAdmin)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "Unable to get Users")
		return
	}
	utils.RespondJSON(w, http.StatusCreated, struct {
		Message   string        `json:"message"`
		SubAdmins []models.User `json:"subAdmins"`
	}{
		Message:   "Get subAdmin successfully.",
		SubAdmins: subAdmins,
	})
}

func RegisterAdmin() {
	logrus.Printf("Creating Admin")

	exists, existsErr := dbHelper.IsUserRoleExists(os.Getenv("ADMIN_EMAIL"), models.RoleAdmin)
	if existsErr != nil {
		logrus.Printf("User Exist: %s", existsErr)
		return
	}
	if exists {
		logrus.Printf("User Exist.")
		return
	}
	hashedPassword, hasErr := utils.HashPassword(os.Getenv("ADMIN_FIRST_PASSWORD"))
	if hasErr != nil {
		logrus.Printf("Unable to Hash Password: %s", existsErr)
		return
	}
	txErr := database.Tx(func(tx *sqlx.Tx) error {
		userID, userExistsErr := dbHelper.IsUserExists(os.Getenv("ADMIN_EMAIL"))
		if userExistsErr != nil {
			logrus.Printf("User Exist: %s", userExistsErr)
			return userExistsErr
		}
		if len(userID) > 0 {
			roleErr := dbHelper.CreateUserRole(tx, userID, userID, models.RoleAdmin)
			if roleErr != nil {
				logrus.Printf("User Role: %s", roleErr)
				return roleErr
			}
		} else {
			userID, saveErr := dbHelper.CreateUser(tx, os.Getenv("ADMIN_NAME"), os.Getenv("ADMIN_EMAIL"), hashedPassword)
			if saveErr != nil {
				logrus.Printf("User Save: %s", saveErr)
				return saveErr
			}
			roleErr := dbHelper.CreateUserRole(tx, userID, userID, models.RoleAdmin)
			if roleErr != nil {
				logrus.Printf("User My Role: %s", roleErr)
				return roleErr
			}
		}
		return nil
	})
	if txErr != nil {
		return
	}
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	Users, err := dbHelper.GetAllUsers(models.RoleUser)
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

func GetAdminUsers(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "subAdmin")
	if id == "" {
		utils.RespondError(w, http.StatusInternalServerError, nil, "Invalid Admin")
	}
	Users, err := dbHelper.GetUsers(id, models.RoleUser)
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

func RemoveSubAdmin(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "userId")
	adminCtx := middlewares.UserContext(r)
	err := dbHelper.RemoveMyUser(id, adminCtx.ID, models.RoleSubAdmin)
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

func RemoveUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "userId")
	err := dbHelper.RemoveUser(id, models.RoleSubAdmin)
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

// Restaurants

func GetAllRestaurants(w http.ResponseWriter, r *http.Request) {
	isFilter := r.URL.Query().Get("isFilter")
	Name := r.URL.Query().Get("name")
	SortBy := r.URL.Query().Get("sortBy")
	Range := r.URL.Query().Get("range")
	if isFilter == "true" {
		var NumberRange int64
		if Name == "" {
			Name = "id"
		}
		if SortBy == "" {
			SortBy = "CreatedAt"
		}
		if Range == "" {
			NumberRange = 10
		} else {
			Number, err := strconv.ParseInt(Range, 10, 64)
			if err != nil {
				utils.RespondError(w, http.StatusInternalServerError, err, "Range is not a number")
				return
			}
			NumberRange = Number
		}
		Restaurants, err := dbHelper.GetFilteredRestaurants(Name, SortBy, NumberRange)
		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, err, "Unable to get Restaurants")
			return
		}
		utils.RespondJSON(w, http.StatusCreated, struct {
			Message     string              `json:"message"`
			Restaurants []models.Restaurant `json:"restaurants"`
		}{
			Message:     "Get Restaurants successfully.",
			Restaurants: Restaurants,
		})

	} else {
		Restaurants, err := dbHelper.GetAllRestaurants()
		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, err, "Unable to get Restaurants")
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
}

func GetSubAdminRestaurants(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "subAdminID")
	Restaurants, err := dbHelper.GetRestaurantsByUserID(id)
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

func GetAllDishes(w http.ResponseWriter, r *http.Request) {
	isFilter := r.URL.Query().Get("isFilter")
	Name := r.URL.Query().Get("name")
	SortBy := r.URL.Query().Get("sortBy")
	Range := r.URL.Query().Get("range")
	if isFilter == "true" {
		var NumberRange int64
		if Name == "" {
			Name = "id"
		}
		if SortBy == "" {
			SortBy = "CreatedAt"
		}
		if Range == "" {
			NumberRange = 10
		} else {
			Number, err := strconv.ParseInt(Range, 10, 64)
			if err != nil {
				utils.RespondError(w, http.StatusInternalServerError, err, "Range is not a number")
				return
			}
			NumberRange = Number
		}
		Dishes, err := dbHelper.GetFilteredDishes(Name, SortBy, NumberRange)
		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, err, "Unable to get Users")
			return
		}
		utils.RespondJSON(w, http.StatusCreated, struct {
			Message string          `json:"message"`
			Dishes  []models.Dishes `json:"dishes"`
		}{
			Message: "Get Dishes successfully.",
			Dishes:  Dishes,
		})

	} else {
		Dishes, err := dbHelper.GetAllDishes()
		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, err, "Unable to get Users")
			return
		}
		utils.RespondJSON(w, http.StatusCreated, struct {
			Message string          `json:"message"`
			Dishes  []models.Dishes `json:"dishes"`
		}{
			Message: "Get Dishes successfully.",
			Dishes:  Dishes,
		})
	}
}

func GetSubAdminDishes(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "subAdminID")
	Menu, err := dbHelper.GetDishesByUserID(id)
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

func UpdateRestaurant(w http.ResponseWriter, r *http.Request) {
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

	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.RespondError(w, http.StatusBadRequest, parseErr, "failed to parse request body")
		return
	}

	restaurant, restaurantErr := dbHelper.GetRestaurantByID(body.ID)
	if restaurantErr != nil {
		utils.RespondError(w, http.StatusBadRequest, nil, "Restaurant not exist")
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

func CloseRestaurant(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "restaurantId")
	err := dbHelper.CloseRestaurant(id)
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

func UpdateDish(w http.ResponseWriter, r *http.Request) {
	body := struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Quantity    int64  `json:"quantity"`
		Price       int64  `json:"price"`
		Discount    int64  `json:"discount"`
	}{}

	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.RespondError(w, http.StatusBadRequest, parseErr, "failed to parse request body")
		return
	}

	dish, dishErr := dbHelper.GetDishByID(body.ID)
	if dishErr != nil {
		utils.RespondError(w, http.StatusBadRequest, nil, "Dish not exist")
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

func RemoveDish(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "dishId")
	err := dbHelper.RemoveDish(id)
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
