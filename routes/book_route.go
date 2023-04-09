package routes

import (
	"github.com/KevDev99/dermatologie24-go-api/controllers"
	"github.com/KevDev99/dermatologie24-go-api/middleware"
	"github.com/gorilla/mux"
)

func BookingRoute(router *mux.Router) {
	router.Handle("/booking", middleware.AuthMiddleware(controllers.AddBooking())).Methods("POST")
	router.Handle("/booking/{id}", middleware.AuthMiddleware(controllers.GetBooking())).Methods("GET")
	router.Handle("/booking/{id}", middleware.AuthMiddleware(controllers.DeleteBooking())).Methods("DELETE")
	router.Handle("/booking/{id}", middleware.AuthMiddleware(controllers.UpdateBooking())).Methods("PATCH")

	router.HandleFunc("/booking/{id}/add-file", middleware.AuthMiddleware(controllers.AddFileToBooking())).Methods("POST")
	router.HandleFunc("/booking/{id}/delete-file/{bookingFileId}", middleware.AuthMiddleware(controllers.DeleteBookingFile())).Methods("DELETE")
}
