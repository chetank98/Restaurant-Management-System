package handler

import (
	"net/http"
	"os"
	"rms/database"
	"rms/database/dbHelper"
	"rms/middlewares"
	"rms/models"
	"rms/utils"

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
