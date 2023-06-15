package main

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

func (app *application) routes() http.Handler {

	router := chi.NewRouter()

	router.Group(func(router chi.Router) {
		router.Post("/user/register", app.registerUserHandler)
		router.Post("/user/auth", app.authorizeUserHandler)
	})

	router.Group(func(router chi.Router) {
		router.Use(app.AuthMiddleware)

		router.Get("/user/{name}", app.getUserByNameHandler)
		router.Post("/user/phone", app.CreateUserPhoneNumberHandler)
		router.Get("/user/phone", app.GetUserPhoneNumberHandler)
		router.Put("/user/phone", app.UpdateUserPhoneNumberHandler)
		router.Delete("/user/phone/{id}", app.DeleteUserPhoneNumberHandler)
	})

	return router
}
