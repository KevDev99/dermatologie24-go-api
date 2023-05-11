package routes

import (
	"github.com/KevDev99/dermatologie24-go-api/controllers"
	"github.com/KevDev99/dermatologie24-go-api/middleware"
	"github.com/gorilla/mux"
)

func StatsRoute(router *mux.Router) {
	router.Handle("/admin/stats", middleware.AuthMiddleware(middleware.AdminMiddleware(controllers.Stats()))).Methods("GET")
}
