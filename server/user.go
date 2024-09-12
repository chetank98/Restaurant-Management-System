package server

import (
	"rms/handler"

	"github.com/go-chi/chi/v5"
)

func adminRoutes(r chi.Router) {
	r.Group(func(admin chi.Router) {
		admin.Post("/subAdmin", handler.RegisterSubAdmin)
		admin.Get("/subAdmins", handler.GetSubAdmins)
		admin.Get("/allUsers", handler.GetAllUsers)
		admin.Get("/allUsers/{id}", handler.GetAdminUsers)
	})
}

func subAdminRoutes(r chi.Router) {
	r.Group(func(subAdmin chi.Router) {
		subAdmin.Get("/", handler.GetInfo)
		subAdmin.Delete("/logout", handler.Logout)
		subAdmin.Post("/user", handler.RegisterUser)
		subAdmin.Get("/users", handler.GetUsers)
		subAdmin.Put("/", handler.UpdateSelfInfo)
	})
}

func userRoutes(r chi.Router) {
	r.Group(func(user chi.Router) {
		user.Get("/", handler.GetInfo)
		user.Delete("/logout", handler.Logout)
		user.Put("/", handler.UpdateSelfInfo)
		user.Post("/address", handler.AddAddress)
		user.Put("/address", handler.UpdateAddress)
	})
}
