package handler

import (
	"net/http"
	"rms/database"
	"rms/database/dbHelper"
	"rms/middlewares"
	"rms/models"
	"rms/utils"

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
