package routes

import (
	"github.com/KevDev99/dermatologie24-go-api/controllers"
	"github.com/KevDev99/dermatologie24-go-api/middleware"
	"github.com/gorilla/mux"
)

func BookingRoute(router *mux.Router) {
	router.Handle("/order", middleware.AuthMiddleware(controllers.AddOrder())).Methods("POST")
	router.Handle("/order/{id}", middleware.AuthMiddleware(controllers.GetOrder())).Methods("GET")
	router.Handle("/order/{id}", middleware.AuthMiddleware(controllers.DeleteOrder())).Methods("DELETE")
	router.Handle("/order/{id}", middleware.AuthMiddleware(controllers.UpdateOrder())).Methods("PATCH")

	router.HandleFunc("/order/{id}/add-file", middleware.AuthMiddleware(controllers.AddFileToOrder())).Methods("POST")
	router.HandleFunc("/order/{id}/delete-file/{bookingFileId}", middleware.AuthMiddleware(controllers.DeleteOrderFile())).Methods("DELETE")
}
