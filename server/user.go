package server

import (
	"rms/handler"

	"github.com/go-chi/chi/v5"
)

func adminRoutes(r chi.Router) {
	r.Group(func(admin chi.Router) {
		admin.Post("/subAdmin", handler.RegisterSubAdmin)
		admin.Get("/subAdmins", handler.GetSubAdmins)
		admin.Delete("/subAdmin/{subAdminId}", handler.RemoveSubAdmin)
	})
}

func subAdminRoutes(r chi.Router) {
	r.Group(func(subAdmin chi.Router) {
		subAdmin.Post("/user", handler.RegisterUser)
		subAdmin.Get("/users", handler.GetUsers)
		subAdmin.Delete("/user/{userId}", handler.RemoveUser)
		subAdmin.Post("/restaurant", handler.OpenRestaurant)
		subAdmin.Put("/restaurant/{restaurantId}", handler.UpdateRestaurant)
		subAdmin.Delete("/restaurant/{restaurantId}", handler.CloseRestaurant)
		subAdmin.Post("/restaurant/{restaurantId}/dish", handler.AddRestaurantDish)
		subAdmin.Put("/restaurant/{restaurantId}/dish/{dishId}", handler.UpdateDish)
		subAdmin.Delete("/restaurant/{restaurantId}/dish/{dishId}", handler.RemoveDish)
	})
}

func userRoutes(r chi.Router) {
	r.Group(func(user chi.Router) {
		user.Post("/address", handler.AddAddress)
		user.Put("/address/{addressId}", handler.UpdateAddress)
		user.Get("/restaurantDistance", handler.GetRestaurantDistance)
	})
}
