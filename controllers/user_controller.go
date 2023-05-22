package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/KevDev99/dermatologie24-go-api/configs"
	"github.com/KevDev99/dermatologie24-go-api/models"
	"github.com/KevDev99/dermatologie24-go-api/utils"
	"github.com/gorilla/context"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func GetUser() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		userId := context.Get(r, "userId")

		var user models.User

		// look up user
		err := configs.DB.Select("id, email, got_patient_details_yn").Where("id = ?", userId).First(&user).Error
		user.Password = ""

		if err != nil {
			utils.SendResponse(rw, http.StatusBadRequest, "error", map[string]interface{}{"data": err.Error()})
			return
		}

		utils.SendResponse(rw, http.StatusOK, "success", map[string]interface{}{"data": user})
	}
}

func GetUserRecipes() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		userId := context.Get(r, "userId")

		var userReceipts []models.UserRecipe

		// look up user
		err := configs.DB.Where("user_id = ?", userId).Find(&userReceipts).Error

		if err != nil {
			utils.SendResponse(rw, http.StatusInternalServerError, "Internal Error.", map[string]interface{}{"data": err.Error})
			return
		}

		utils.SendResponse(rw, http.StatusOK, "success", map[string]interface{}{"data": userReceipts})
	}
}

func DeleteUser() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		userId := context.Get(r, "userId")

		var user models.User

		// look up user
		err := configs.DB.Delete(&user, userId).Error

		if err != nil {
			utils.SendResponse(rw, http.StatusBadRequest, "error", map[string]interface{}{"data": err.Error()})
			return
		}

		utils.SendResponse(rw, http.StatusOK, "error", map[string]interface{}{"data": "user deleted."})
	}
}

func UpdateUser() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		userId := context.Get(r, "userId")

		var user models.User

		// get user and check if user exists
		err := configs.DB.First(&user, userId).Error
		if err != nil {
			utils.SendResponse(rw, http.StatusBadRequest, "error", map[string]interface{}{"data": err.Error()})
			return
		}

		// Decode the request body into a User struct
		var updateUser models.User
		err = json.NewDecoder(r.Body).Decode(&updateUser)
		if err != nil {
			utils.SendResponse(rw, http.StatusBadRequest, "error", map[string]interface{}{"data": err.Error()})
			return
		}

		// update user in DB
		dberr := configs.DB.Model(&user).Updates(updateUser).Error

		// show internal error if needed
		if dberr != nil {
			utils.SendResponse(rw, http.StatusBadRequest, "error", map[string]interface{}{"data": err.Error()})
			return
		}

		// return updates user and successfull status code
		utils.SendResponse(rw, http.StatusOK, "success", map[string]interface{}{"data": user})
	}
}

func GetUsers() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		var users []models.User

		// retrieve all users
		// exclude admin
		err := configs.DB.Preload("UserDetails").Where("adminYN = ?", "0").Find(&users).Error

		// get user and check if user exists
		if err != nil {
			utils.SendResponse(rw, http.StatusInternalServerError, "Internal Error", map[string]interface{}{"data": err.Error()})
			return
		}

		// return updates user and successfull status code
		utils.SendResponse(rw, http.StatusOK, "success", map[string]interface{}{"data": users})
	}
}
