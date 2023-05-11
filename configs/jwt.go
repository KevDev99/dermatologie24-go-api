package configs

import (
	"time"

	"github.com/KevDev99/dermatologie24-go-api/models"
	"github.com/dgrijalva/jwt-go"
)

var JwtKey = []byte("your_secret_key")
var RefreshJwtKey = []byte("your_secret_refresh_key")

type TokenType byte

const (
	AccessToken  TokenType = 0
	RefreshToken TokenType = 1
)

func GenerateToken(user models.User, typ byte) (string, error) {

	var tokenClaims jwt.MapClaims

	if typ == byte(AccessToken) {
		tokenClaims = jwt.MapClaims{
			"id":    user.Id,
			"email": user.Email,
			"exp":   time.Now().Add(time.Hour * 2).Unix(),
		}
	} else {
		tokenClaims = jwt.MapClaims{
			"id":  user.Id,
			"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
		}
	}
	// Create a new token object
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)

	// Generate a signed token string
	var tokenString string
	var err error

	if typ == byte(AccessToken) {
		tokenString, err = token.SignedString(JwtKey)
	} else {
		tokenString, err = token.SignedString(RefreshJwtKey)
	}

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
