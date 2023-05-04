package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/KevDev99/dermatologie24-go-api/configs"
	"github.com/KevDev99/dermatologie24-go-api/models"
	"github.com/KevDev99/dermatologie24-go-api/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/gorilla/context"
)

func Login() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		var user models.User
		var userInput models.User

		// Parse the request body
		err := json.NewDecoder(r.Body).Decode(&userInput)

		if err != nil {
			utils.SendResponse(rw, http.StatusBadRequest, err.Error(), map[string]interface{}{"data": err.Error()})
			return
		}

		// query user
		queryErr := configs.DB.Where("email = ?", userInput.Email).First(&user).Error
		if queryErr != nil {
			utils.SendResponse(rw, http.StatusBadRequest, "Wrong Credentials. Please try again.", map[string]interface{}{"data": queryErr.Error()})
			return
		}

		// compare pw
		authErr := utils.ComparePasswords([]byte(user.Password), userInput.Password)
		if authErr != nil {
			utils.SendResponse(rw, http.StatusBadRequest, "Wrong Credentials. Please try again.", map[string]interface{}{"data": authErr.Error()})
			return
		}

		// Generate the JWT token
		accessToken, accessTokenErr := configs.GenerateToken(user, byte(configs.AccessToken))

		if accessTokenErr != nil {
			utils.SendResponse(rw, http.StatusInternalServerError, "error", map[string]interface{}{"data": accessTokenErr.Error()})
			return
		}

		refreshToken, refreshTokenErr := configs.GenerateToken(user, byte(configs.RefreshToken))

		if refreshTokenErr != nil {
			utils.SendResponse(rw, http.StatusInternalServerError, "error", map[string]interface{}{"data": refreshTokenErr.Error()})
			return
		}

		// Set the JWT token in the response header
		rw.Header().Set("Authorization", "Bearer "+accessToken)

		// Return a success response
		utils.SendResponse(rw, http.StatusOK, "success", map[string]interface{}{"data": map[string]interface{}{"access_token": accessToken, "refresh_token": refreshToken}})
	}
}

func Protected() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		userId := context.Get(r, "userId")
		fmt.Fprintln(rw, "your user id is:", userId)
	}
}

func Register() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var user models.User

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			utils.SendResponse(rw, http.StatusBadRequest, "error", map[string]interface{}{"data": err.Error()})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&user); validationErr != nil {
			utils.SendResponse(rw, http.StatusBadRequest, "error", map[string]interface{}{"data": validationErr.Error()})
			return
		}

		// hash pw
		hashedPassword, hashError := utils.GetHashedPassword(user.Password)
		if hashError != nil {
			utils.SendResponse(rw, http.StatusInternalServerError, "error", map[string]interface{}{"data": hashError.Error()})
			return
		}

		// create new user
		newUser := models.User{
			Firstname:           user.Firstname,
			Lastname:            user.Lastname,
			Email:               user.Email,
			Password:            string(hashedPassword),
			GotPatientDetailsYN: false,
		}

		if err := configs.DB.Create(&newUser).Error; err != nil {
			utils.SendResponse(rw, http.StatusBadRequest, "error", map[string]interface{}{"data": err.Error()})
			return
		}

		utils.SendResponse(rw, http.StatusOK, "success", map[string]interface{}{"data": newUser})
	}
}

func RefreshToken() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		var tokenDTO models.Token

		// Get the refresh token from the request
		err := json.NewDecoder(r.Body).Decode(&tokenDTO)

		if err != nil {
			utils.SendResponse(rw, http.StatusBadRequest, "error", map[string]interface{}{"data": err.Error()})
			return
		}

		// Verify that the refresh token is not empty
		if tokenDTO.RefreshToken == "" {
			utils.SendResponse(rw, http.StatusBadRequest, "error", map[string]interface{}{"data": "refresh token is missing."})
			return
		}

		// auth with token
		token, err := jwt.Parse(tokenDTO.RefreshToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return configs.RefreshJwtKey, nil
		})

		// error handling
		if err != nil || !token.Valid {
			utils.SendResponse(rw, http.StatusUnauthorized, "error", map[string]interface{}{"data": "key couldnt be verified as  valid."})
			return
		}

		// Extract the user ID from the refresh token
		userId := token.Claims.(jwt.MapClaims)["id"].(float64)

		// Check if the user exists
		var user models.User
		queryErr := configs.DB.Where("id = ?", userId).First(&user).Error
		if queryErr != nil {
			utils.SendResponse(rw, http.StatusBadRequest, "error", map[string]interface{}{"data": queryErr.Error()})
			return
		}

		// Generate a new access token
		accessToken, accessTokenErr := configs.GenerateToken(user, byte(configs.AccessToken))
		if accessTokenErr != nil {
			utils.SendResponse(rw, http.StatusBadRequest, "error", map[string]interface{}{"data": accessTokenErr.Error()})
			return
		}

		// Generate a new refresh token
		refreshToken, refreshTokenErr := configs.GenerateToken(user, byte(configs.RefreshToken))
		if refreshTokenErr != nil {
			utils.SendResponse(rw, http.StatusBadRequest, "error", map[string]interface{}{"data": refreshTokenErr.Error()})
			return
		}

		// Return a success response
		utils.SendResponse(rw, http.StatusOK, "success", map[string]interface{}{"access_token": accessToken, "refresh_token": refreshToken})

	}
}

func PasswordReset() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var user models.User

		// check if token is set -> user has already received token and can now enter a new one
		queryParams := r.URL.Query()
		token := queryParams.Get("token")

		if token != "" {
			// create the custom URL scheme to launch the app
			url := os.Getenv("APP_URL") + "?token=" + token

			// redirect to the custom URL scheme
			http.Redirect(rw, r, url, http.StatusSeeOther)
			return
		}

		// Check if the "field" form field is included in the request body
		email := r.FormValue("email")

		if email == "" {
			utils.SendResponse(rw, http.StatusBadRequest, "No email provided", map[string]interface{}{"data": "No email provided"})
			return
		}

		// get user by mail
		queryErr := configs.DB.Where("email = ?", email).First(&user).Error
		if queryErr != nil {
			utils.SendResponse(rw, http.StatusBadRequest, "No user with that email found.", map[string]interface{}{"data": queryErr.Error()})
			return
		}

		// check if there is already a token in the db
		newToken := uuid.New().String()
		expiresAt := time.Now().Add(time.Hour * 24)

		// create new password forget token
		passwordForgetToken := models.PasswordResetToken{Token: newToken, UserID: user.Id, ExpiresAt: expiresAt}

		// add to db
		err := configs.DB.Create(&passwordForgetToken).Error
		if err != nil {
			utils.SendResponse(rw, http.StatusInternalServerError, "Internal server error", map[string]interface{}{"data": err.Error()})
			return
		}

		htmlBody := fmt.Sprintf("<h1>Passwort Zurücksetzen</h1><br/><br/><h2>Url: <a href='www.google.com'>http://localhost:5000/reset-password?token=%s</a> <h2>", passwordForgetToken.Token)

		configs.SendMail("no-reply@dermatologie24.com", "kevin.taufer@outlook.com", "Password zurücksetzen", htmlBody)

		utils.SendResponse(rw, http.StatusOK, "reset token send via mail.", map[string]interface{}{"data": "reset token send via mail."})
	}
}
