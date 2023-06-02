package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/KevDev99/dermatologie24-go-api/configs"
	"github.com/KevDev99/dermatologie24-go-api/models"
	"github.com/KevDev99/dermatologie24-go-api/responses"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func GetHashedPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	return hashedPassword, err
}

func ComparePasswords(hashedPw []byte, passwordToCompare string) error {
	return bcrypt.CompareHashAndPassword(hashedPw, []byte(passwordToCompare))
}

func Contains(stringSlice []string, text string) bool {
	for _, a := range stringSlice {
		if a == text {
			return true
		}
	}
	return false
}

func SendResponse(rw http.ResponseWriter, httpStatusCode int, message string, responseData map[string]interface{}) {
	rw.WriteHeader(httpStatusCode)
	response := responses.GeneralResponse{Status: httpStatusCode, Message: message, Data: responseData}
	json.NewEncoder(rw).Encode(response)
}

func CheckIfTokenNotExpired(endTime time.Time, tokenTime time.Time) bool {

	return endTime.After(tokenTime)
}

func SendMailConfirmationMail(userId int, email string) error {
	// create token for email confirmation
	// check if there is already a token in the db
	newToken := uuid.New().String()
	expiresAt := time.Now().Add(time.Hour * 24)
	emailConfirmToken := models.EmailConfirmToken{Token: newToken, UserID: userId, ExpiresAt: expiresAt}

	err := configs.DB.Create(&emailConfirmToken).Error

	appUrl := os.Getenv("APP_URL") + fmt.Sprintf("/email-confirm?token=%s", emailConfirmToken.Token)

	payload := fmt.Sprintf(`{
		"sender": {
			"email": "kevin.taufer@outlook.com",
			"name": "kevin.taufer@outlook.com"
		},
		"params": {
					"link": "%s"
				},
		"templateId": %s,
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
				"subject": "Email Bestätigen"
			}
		]
	}`, appUrl, os.Getenv("BREVO_MAIL_CONFIRM_TEMPLATE_ID"), email, appUrl)

	configs.SendMail(payload)

	return err
}
