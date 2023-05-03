package routes

import (
	"github.com/KevDev99/dermatologie24-go-api/controllers"
	"github.com/KevDev99/dermatologie24-go-api/middleware"
	"github.com/gorilla/mux"
)

func StatsRoute(router *mux.Router) {
	router.Handle("/stats", middleware.AuthMiddleware(controllers.Stats())).Methods("GET")
}
