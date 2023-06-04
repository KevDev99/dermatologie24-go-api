package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/KevDev99/dermatologie24-go-api/configs"
	"github.com/KevDev99/dermatologie24-go-api/models"
	"github.com/KevDev99/dermatologie24-go-api/utils"
	"github.com/gorilla/mux"
)

func GetRecipes() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// userId := context.Get(r, "userId")

		var recipes []models.Recipe

		// look up user
		err := configs.DB.Find(&recipes).Error

		if err != nil {
			utils.SendResponse(rw, http.StatusInternalServerError, "Internal Error.", map[string]interface{}{"data": err.Error})
			return
		}

		utils.SendResponse(rw, http.StatusOK, "success", map[string]interface{}{"data": recipes})
	}
}

func GetRecipe() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		recipeId := mux.Vars(r)["id"]

		var recipe models.Recipe

		// look up user
		err := configs.DB.Where("id = ?", recipeId).First(&recipe).Error

		if err != nil {
			utils.SendResponse(rw, http.StatusInternalServerError, "Internal Error.", map[string]interface{}{"data": err.Error})
			return
		}

		utils.SendResponse(rw, http.StatusOK, "success", map[string]interface{}{"data": recipe})
	}
}

func DeleteRecipe() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		recipeId := mux.Vars(r)["id"]

		var recipe models.Recipe

		// delete recipe
		err := configs.DB.Delete(&recipe, recipeId).Error

		if err != nil {
			utils.SendResponse(rw, http.StatusInternalServerError, "Internal Error.", map[string]interface{}{"data": err.Error})
			return
		}

		utils.SendResponse(rw, http.StatusOK, "success", map[string]interface{}{"data": "successfully deleted."})
	}
}

func AddRecipe() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var recipe models.Recipe

		// Parse the incoming request as a multipart form with a maximum file size of 32MB
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			utils.SendResponse(rw, http.StatusBadRequest, "error", map[string]interface{}{"data": err.Error()})
			return
		}

		files := r.MultipartForm.File["files"] // Get the list of uploaded files from the request
		form := r.Form

		recipe.Title = form.Get("title")
		recipe.Information = form.Get("information")

		// create order
		configs.DB.Create(&recipe)

		// upload files
		go uploadFiles(files)

		utils.SendResponse(rw, http.StatusCreated, "success", map[string]interface{}{"data": recipe})
	}
}

func UpdateRecipe() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		recipeId := mux.Vars(r)["id"]

		var recipe models.Recipe

		err := configs.DB.First(&recipe, recipeId).Error
		if err != nil {
			utils.SendResponse(rw, http.StatusBadRequest, "Recipe Not Found.", map[string]interface{}{"data": err.Error()})
			return
		}

		// decode the request body into a User struct
		var updateRecipe models.Recipe
		err = json.NewDecoder(r.Body).Decode(&updateRecipe)
		if err != nil {
			utils.SendResponse(rw, http.StatusBadRequest, "Internal Error.", map[string]interface{}{"data": err.Error()})
			return
		}

		// update  in DB
		dberr := configs.DB.Model(&recipe).Updates(updateRecipe).Error

		// show internal error if needed
		if dberr != nil {
			utils.SendResponse(rw, http.StatusInternalServerError, "Internal Error", map[string]interface{}{"data": dberr.Error()})
			return
		}

		utils.SendResponse(rw, http.StatusOK, "success", map[string]interface{}{"data": recipe})
	}
}
