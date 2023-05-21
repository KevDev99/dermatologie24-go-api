package routes

import (
	"github.com/KevDev99/dermatologie24-go-api/controllers"
	"github.com/KevDev99/dermatologie24-go-api/middleware"
	"github.com/gorilla/mux"
)

func AuthRoute(router *mux.Router) {
	router.Handle("/auth/login", controllers.Login()).Methods("POST")
	router.Handle("/auth/protected", middleware.AuthMiddleware(controllers.Protected())).Methods("POST")
	router.Handle("/auth/register", controllers.Register()).Methods("POST")
	router.Handle("/auth/refresh-token", controllers.RefreshToken()).Methods("POST")
	router.Handle("/admin/login", controllers.AdminLogin()).Methods("POST")
	router.Handle("/update-password", middleware.AuthMiddleware(controllers.UpdatePassword())).Methods("POST")
	router.Handle("/update-email", middleware.AuthMiddleware(controllers.UpdateEmail())).Methods("PATCH")
}
