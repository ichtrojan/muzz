package routes

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/ichtrojan/muzz/controllers"
	"github.com/ichtrojan/muzz/helpers"
	customMiddleware "github.com/ichtrojan/muzz/middleware"
	"net/http"
)

func ApiRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "X-api-key", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Use(middleware.Logger)

	// force `content-type` to `application/json`
	router.Use(customMiddleware.AcceptJson)

	// response for 404
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(helpers.PrepareMessage("endpoint not found"))
		return
	})

	// Login
	router.Post("/login", controllers.Login)

	router.Post("/user/create", controllers.CreateUser)

	router.Group(func(router chi.Router) {
		// auth middleware
		router.Use(customMiddleware.AuthenticateUser)

		// discover
		router.Get("/discover", controllers.Discover)

		// swipe user
		router.Post("/swipe", controllers.Swipe)
	})

	return router
}
