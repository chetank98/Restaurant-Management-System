package handler

import (
	"net/http"
	"rms/database"
	"rms/database/dbHelper"
	"rms/middlewares"
	"rms/models"
	"rms/utils"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func LoginUser(w http.ResponseWriter, r *http.Request) {
	body := struct {
		Email    string      `json:"email"`
		Password string      `json:"password"`
		Role     models.Role `json:"role"`
	}{}

	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.RespondError(w, http.StatusBadRequest, parseErr, "failed to parse request body")
		return
	}

	userId, userRoleId, userErr := dbHelper.GetUserRoleIDByPassword(body.Email, body.Password, body.Role)
	if userErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, userErr, "failed to find user")
		return
	}
	// create user session
	sessionToken, jwtError := utils.JwtToken(userId, userRoleId)
	if jwtError != nil {
		utils.RespondError(w, http.StatusInternalServerError, jwtError, jwtError.Error())
		return
	}
	sessionErr := dbHelper.CreateUserSession(database.RMS, userId, userRoleId, sessionToken)
	if sessionErr != nil {
		utils.RespondError(w, http.StatusInternalServerError, sessionErr, "failed to create user session")
		return
	}
	utils.RespondJSON(w, http.StatusCreated, struct {
		Token string `json:"token"`
		Type  string `json:"tokenType"`
	}{
		Token: sessionToken,
		Type:  "Bearer",
	})
}

func GetInfo(w http.ResponseWriter, r *http.Request) {
	userCtx := middlewares.UserContext(r)
	utils.RespondJSON(w, http.StatusOK, userCtx)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	token := strings.Split(r.Header.Get("authorization"), " ")[1]
	err := dbHelper.DeleteSessionToken(token)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed to logout user")
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func UpdateSelfInfo(w http.ResponseWriter, r *http.Request) {
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
	if body.Name == "" {
		body.Name = adminCtx.Name
	}
	if !utils.IsEmailValid(body.Email) {
		if body.Email != "" {
			utils.RespondError(w, http.StatusBadRequest, nil, "fInvalid Email")
			return
		}
		body.Email = adminCtx.Email
	}
	if body.Password == "" {
		body.Password = adminCtx.Password
	} else {
		hashedPassword, hasErr := utils.HashPassword(body.Password)
		if hasErr != nil {
			utils.RespondError(w, http.StatusInternalServerError, hasErr, "failed to secure password")
			return
		}
		body.Password = hashedPassword
	}
	err := dbHelper.UpdateUserInfo(adminCtx.ID, body.Name, body.Email, body.Password)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "failed update User")
		return
	}
	utils.RespondJSON(w, http.StatusCreated, struct {
		Message string `json:"message"`
	}{
		Message: "user Update successfully",
	})
}

func AddAddress(w http.ResponseWriter, r *http.Request) {
	body := struct {
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
		addressErr := dbHelper.CreateUserAddress(tx, userCtx.ID, body.Address, body.State, body.City, body.PinCode, body.Lat, body.Lng)
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
		Message: "Address Created successfully",
	})
}

func UpdateAddress(w http.ResponseWriter, r *http.Request) {
	body := struct {
		ID      string  `json:"addressId"`
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

	userAddress, addressErr := utils.GetUserAddressById(body.ID, userCtx.UserAddresses)
	if addressErr != nil {
		utils.RespondError(w, http.StatusBadRequest, nil, "Address not exist")
		return
	}

	if len(body.Address) > 30 || len(body.Address) <= 2 {
		if body.Address != "" {
			utils.RespondError(w, http.StatusBadRequest, nil, "Address must be with in 2 to 30 letter.")
			return
		}
		body.Address = userAddress.Address
	}

	if len(body.State) > 16 || len(body.State) <= 2 {
		if body.State != "" {
			utils.RespondError(w, http.StatusBadRequest, nil, "State must be with in 2 to 16 letter.")
			return
		}
		body.State = userAddress.State
	}

	if len(body.City) > 20 || len(body.City) <= 2 {
		if body.City != "" {
			utils.RespondError(w, http.StatusBadRequest, nil, "City must be with in 2 to 20 letter.")
			return
		}
		body.City = userAddress.City
	}

	if len(body.PinCode) != 6 {
		if body.PinCode != "" {
			utils.RespondError(w, http.StatusBadRequest, nil, "PinCode must 6 digit.")
			return
		}
		body.PinCode = userAddress.PinCode
	}

	if body.Lat > 90 || body.Lat < -90 {
		if body.Lat != 0 {
			utils.RespondError(w, http.StatusBadRequest, nil, "Invalid Latitude.")
			return
		}
		body.Lat = userAddress.Lat
	}

	if body.Lng > 180 || body.Lng < -180 {
		if body.Lat != 0 {
			utils.RespondError(w, http.StatusBadRequest, nil, "Invalid Longitude.")
			return
		}
		body.Lng = userAddress.Lng
	}

	err := dbHelper.UpdateUserAddress(body.ID, body.Address, body.State, body.City, body.PinCode, body.Lat, body.Lng)
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

func GetRestaurantsDish(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "restaurantId")
	RestaurantsMenu, err := dbHelper.GetRestaurantDish(id)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "Unable to get Users")
		return
	}
	utils.RespondJSON(w, http.StatusCreated, struct {
		Message         string          `json:"message"`
		RestaurantsMenu []models.Dishes `json:"restaurantsMenu"`
	}{
		Message:         "Get Restaurants successfully.",
		RestaurantsMenu: RestaurantsMenu,
	})
}

func GetRestaurantsDistance(w http.ResponseWriter, r *http.Request) {
	restaurantId := r.URL.Query().Get("restaurantId")
	userCtx := middlewares.UserContext(r)
	addressId := r.URL.Query().Get("addressId")

	Restaurant, err := dbHelper.GetRestaurantByID(restaurantId)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, "Unable to get Restaurant")
		return
	}

	userAddress, addressErr := utils.GetUserAddressById(addressId, userCtx.UserAddresses)
	if addressErr != nil {
		utils.RespondError(w, http.StatusBadRequest, nil, "Address not exist")
		return
	}

	Distance, Unit := utils.CalculateDistance(userAddress.Lat, userAddress.Lng, Restaurant.Lat, Restaurant.Lng)
	utils.RespondJSON(w, http.StatusCreated, struct {
		Message      string  `json:"message"`
		Distance     float64 `json:"restaurantsDistance"`
		DistanceUnit string  `json:"distanceUnit"`
	}{
		Message:      "Get Restaurants successfully.",
		Distance:     Distance,
		DistanceUnit: Unit,
	})
}
