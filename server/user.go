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
		admin.Get("/allUsers/{subAdmin}", handler.GetAdminUsers)
		admin.Get("/allRestaurants", handler.GetAllRestaurants)
		admin.Get("/allDishes", handler.GetAllDishes)
		admin.Get("/subAdminRestaurants/{subAdminID}", handler.GetSubAdminRestaurants)
		admin.Get("/subAdminDishes/{subAdminID}", handler.GetSubAdminDishes)
		admin.Put("/restaurant", handler.UpdateRestaurant)
		admin.Put("/dish", handler.UpdateRestaurant)
	})
}

func subAdminRoutes(r chi.Router) {
	r.Group(func(subAdmin chi.Router) {
		subAdmin.Get("/", handler.GetInfo)
		subAdmin.Delete("/logout", handler.Logout)
		subAdmin.Post("/user", handler.RegisterUser)
		subAdmin.Get("/users", handler.GetUsers)
		subAdmin.Put("/", handler.UpdateSelfInfo)
		subAdmin.Post("/restaurant", handler.OpenRestaurant)
		subAdmin.Post("/restaurantsDish", handler.AddRestaurantDish)
		subAdmin.Get("/myRestaurants", handler.GetMyRestaurants)
		subAdmin.Get("/myDishes", handler.GetMyDishes)
		subAdmin.Get("/restaurantDish/{restaurantId}", handler.GetRestaurantsDish)
		subAdmin.Put("/myRestaurants", handler.UpdateMyRestaurant)
		subAdmin.Put("/myDish", handler.UpdateMyRestaurant)
	})
}

func userRoutes(r chi.Router) {
	r.Group(func(user chi.Router) {
		user.Get("/", handler.GetInfo)
		user.Delete("/logout", handler.Logout)
		user.Put("/", handler.UpdateSelfInfo)
		user.Post("/address", handler.AddAddress)
		user.Put("/address", handler.UpdateAddress)
		user.Get("/allRestaurants", handler.GetAllRestaurants)
		user.Get("/restaurantDish/{restaurantId}", handler.GetRestaurantsDish)
		user.Get("/restaurantDistance", handler.GetRestaurantsDistance)
	})
}
