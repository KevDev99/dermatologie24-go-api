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
	router.Handle("/user-recipes", middleware.AuthMiddleware(controllers.GetUserRecipes())).Methods("GET")

	router.Handle("/reset-password", controllers.PasswordReset()).Methods("POST")
	router.Handle("/confirm-mail", controllers.EmailConfirm()).Methods("POST")
	router.Handle("/admin/users", middleware.AuthMiddleware(middleware.AdminMiddleware(controllers.GetUsers()))).Methods("GET")

	router.Handle("/user/profile-data", middleware.AuthMiddleware(controllers.ProfileData())).Methods("POST")
}
