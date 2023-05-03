package routes

import (
	"github.com/KevDev99/dermatologie24-go-api/controllers"
	"github.com/KevDev99/dermatologie24-go-api/middleware"
	"github.com/gorilla/mux"
)

func UserRoute(router *mux.Router) {
	router.Handle("/user", middleware.AuthMiddleware(controllers.GetUser())).Methods("GET")
	router.Handle("/user", middleware.AuthMiddleware(controllers.DeleteUser())).Methods("DELETE")
	router.Handle("/user", middleware.AuthMiddleware(controllers.UpdateUser())).Methods("PATCH")
	router.Handle("/reset-password", controllers.PasswordReset()).Methods("POST")
}
