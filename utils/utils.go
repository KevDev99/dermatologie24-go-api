package utils

import (
	"encoding/json"
	"net/http"

	"github.com/KevDev99/dermatologie24-go-api/responses"
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
