package routes

import (
	"github.com/KevDev99/dermatologie24-go-api/controllers"
	"github.com/KevDev99/dermatologie24-go-api/middleware"
	"github.com/gorilla/mux"
)

func RecipeRoute(router *mux.Router) {
	router.Handle("/recipes", middleware.AdminMiddleware(middleware.AuthMiddleware(controllers.GetRecipes()))).Methods("GET")
	router.Handle("/recipe/{id}", middleware.AdminMiddleware(middleware.AuthMiddleware(controllers.GetRecipe()))).Methods("GET")
	router.Handle("/recipe/{id}", middleware.AdminMiddleware(middleware.AuthMiddleware(controllers.DeleteRecipe()))).Methods("DELETE")
	router.Handle("/recipe", middleware.AdminMiddleware(middleware.AuthMiddleware(controllers.AddRecipe()))).Methods("POST")
}
