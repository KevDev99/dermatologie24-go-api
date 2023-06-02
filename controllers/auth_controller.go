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

		if !user.EmailConfirmedYN {
			utils.SendResponse(rw, http.StatusBadRequest, "Email not confirmed.\nPlease check your inbox.", map[string]interface{}{"data": "Email not confirmed."})
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

		var userExisting models.User

		configs.DB.First(&userExisting, "email = ?", user.Email)
		if userExisting.Email != "" {
			utils.SendResponse(rw, http.StatusBadRequest, "User already existing.", map[string]interface{}{"data": "User already existing."})
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
			Email:               user.Email,
			Password:            string(hashedPassword),
			GotPatientDetailsYN: false,
		}

		if err := configs.DB.Create(&newUser).Error; err != nil {
			utils.SendResponse(rw, http.StatusBadRequest, "error", map[string]interface{}{"data": err.Error()})
			return
		}

		if err := utils.SendMailConfirmationMail(newUser.Id, newUser.Email); err != nil {
			utils.SendResponse(rw, http.StatusBadRequest, "Internal Error", map[string]interface{}{"data": err.Error()})
			return
		}

		// clear pw from response
		newUser.Password = ""

		utils.SendResponse(rw, http.StatusOK, "User created.", map[string]interface{}{"data": newUser})
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

		appUrl := os.Getenv("APP_URL") + fmt.Sprintf("/reset?token=%s", passwordForgetToken.Token)

		payload := fmt.Sprintf(`{
			"sender":{
				"email": "%s",
				"name": "%s"
			},
			"subject": "%s",
			"templateId": %s,
			"params":{
				"link": "%s"
			},
			"messageVersions": [
				{
					"to": [
						{
							"email": "%s"
						}
					],
					"params": {
						"link": "%s"
					},
					"subject": "Passwort zurücksetzen"
				}
			]
		}`, email, email, "Passwort zurücksetzen", os.Getenv("BREVO_PASSWORT_RESET_TEMPLATE_ID"), appUrl, user.Email, appUrl)

		configs.SendMail(payload)

		utils.SendResponse(rw, http.StatusOK, "reset token send via mail.", map[string]interface{}{"data": "reset token send via mail."})
	}
}

func EmailConfirm() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var emailConfirmToken models.EmailConfirmToken

		// check if token is set -> user has already received token and can now enter a new one
		queryParams := r.URL.Query()
		token := queryParams.Get("token")

		if token == "" {
			utils.SendResponse(rw, http.StatusBadRequest, "No Token Provided", map[string]interface{}{"data": "No Token Provided"})
			return
		}

		// get token
		queryErr := configs.DB.Where("token = ?", token).First(&emailConfirmToken).Error
		if queryErr != nil {
			utils.SendResponse(rw, http.StatusBadRequest, "Token Not Found.", map[string]interface{}{"data": "Token Not Found."})
			return
		}

		// validate token

		// set confirmed to true
		configs.DB.Exec("UPDATE users SET email_confirmed_yn = ? WHERE id = ?", 1, emailConfirmToken.UserID)

		// delete from email confirmations table
		configs.DB.Delete(&emailConfirmToken)

		utils.SendResponse(rw, http.StatusOK, "email confirmed", map[string]interface{}{"data": "email confirmed"})
	}
}

func AdminLogin() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		var user models.User
		var userInput models.User

		// Parse the request body
		err := json.NewDecoder(r.Body).Decode(&userInput)

		if err != nil {
			utils.SendResponse(rw, http.StatusBadRequest, err.Error(), map[string]interface{}{"data": err.Error()})
			return
		}

		// query admin
		queryErr := configs.DB.Where("email = ? AND adminYN = ?", userInput.Email, 1).First(&user).Error
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
			utils.SendResponse(rw, http.StatusInternalServerError, "Internal Error", map[string]interface{}{"data": "Internal Error"})
			fmt.Println("Internal Error when creating Access Token: " + accessTokenErr.Error())
			return
		}

		// Set the JWT token in the response header
		rw.Header().Set("Authorization", "Bearer "+accessToken)

		// Return a success response
		utils.SendResponse(rw, http.StatusOK, "success", map[string]interface{}{"data": map[string]interface{}{"access_token": accessToken}})
	}
}

func UpdatePassword() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		type PasswordToUpdate struct {
			Token    string `json:"token" validate:"required"`
			Password string `json:"password" validate:"required"`
		}

		var passwordToUpdate PasswordToUpdate
		var passwordResetToken models.PasswordResetToken

		// Parse the request body
		err := json.NewDecoder(r.Body).Decode(&passwordToUpdate)

		if err != nil {
			utils.SendResponse(rw, http.StatusBadRequest, err.Error(), map[string]interface{}{"data": err.Error()})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&passwordToUpdate); validationErr != nil {
			utils.SendResponse(rw, http.StatusBadRequest, "Token or Password field not provided", map[string]interface{}{"data": validationErr.Error()})
			return
		}

		fmt.Println(passwordToUpdate)

		// get token from database
		queryErr := configs.DB.Where("token = ?", passwordToUpdate.Token).First(&passwordResetToken).Error
		if queryErr != nil {
			utils.SendResponse(rw, http.StatusBadRequest, "Token Not Found.", map[string]interface{}{"data": "Token Not Found."})
			return
		}

		// check if token expired
		if !utils.CheckIfTokenNotExpired(passwordResetToken.ExpiresAt, time.Now()) {
			utils.SendResponse(rw, http.StatusBadRequest, "Token Already Expired.", map[string]interface{}{"data": "Token Already Expired."})
			return
		}

		// hash pw
		hashedPassword, hashError := utils.GetHashedPassword(passwordToUpdate.Password)
		if hashError != nil {
			utils.SendResponse(rw, http.StatusInternalServerError, "Internal Error", map[string]interface{}{"data": hashError.Error()})
			return
		}

		// update password
		configs.DB.Model(&models.User{}).Where("id = ?", passwordResetToken.UserID).Update("password", string(hashedPassword))

		// delete from table
		configs.DB.Delete(&passwordResetToken)

		// Return a success response
		utils.SendResponse(rw, http.StatusOK, "password reseted.", map[string]interface{}{"data": "password reseted."})
	}
}

func UpdateEmail() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		userId := fmt.Sprintf("%d", int(context.Get(r, "userId").(float64)))

		type EmailToUpdate struct {
			NewEmail string `json:"email" validate:"required"`
		}

		var emailToUpdate EmailToUpdate
		var user models.User
		var updateUser models.User

		// Parse the request body
		err := json.NewDecoder(r.Body).Decode(&emailToUpdate)

		if err != nil {
			utils.SendResponse(rw, http.StatusBadRequest, "Invalid Request", map[string]interface{}{"data": err.Error()})
			return
		}

		// get user and check if user exists
		queryErr := configs.DB.First(&user, userId).Error
		if queryErr != nil {
			utils.SendResponse(rw, http.StatusBadRequest, "Der Benutzer konnte nicht gefunden werden.", map[string]interface{}{"data": err.Error()})
			return
		}

		// check if the mail is already selected
		// get token from database
		rowsAffected := configs.DB.Where("email = ?", emailToUpdate.NewEmail).First(&models.User{}).RowsAffected
		if rowsAffected > 0 {
			utils.SendResponse(rw, http.StatusBadRequest, "There is already an user with that email.", map[string]interface{}{"data": "There is already an user with that email."})
			return
		}

		updateUser = user
		updateUser.Email = emailToUpdate.NewEmail
		updateUser.EmailConfirmedYN = false

		// resend email confirm
		if err := utils.SendMailConfirmationMail(updateUser.Id, updateUser.Email); err != nil {
			utils.SendResponse(rw, http.StatusBadRequest, "Internal Error", map[string]interface{}{"data": err.Error()})
			return
		}

		configs.DB.Exec("UPDATE users SET email_confirmed_yn = ? WHERE id = ?", 0, updateUser.Id)

		// update user in DB
		dberr := configs.DB.Model(&user).Updates(updateUser).Error

		// show internal error if needed
		if dberr != nil {
			utils.SendResponse(rw, http.StatusBadRequest, "error", map[string]interface{}{"data": err.Error()})
			return
		}

		utils.SendResponse(rw, http.StatusOK, "success", map[string]interface{}{"data": "success"})
	}
}
