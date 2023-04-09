package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/KevDev99/dermatologie24-go-api/configs"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// check if auth header is present
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		// check if bearer token is set and its valid
		authToken := strings.Split(authHeader, " ")
		if len(authToken) != 2 {
			http.Error(w, "Invalid authorization token", http.StatusUnauthorized)
			return
		}

		// auth with token
		token, err := jwt.Parse(authToken[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return configs.JwtKey, nil
		})

		// error handling
		if err != nil {
			http.Error(w, "Invalid authorization token", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, "Invalid authorization token", http.StatusUnauthorized)
			return
		}

		// Call the next handler with the user ID from the token
		userId := token.Claims.(jwt.MapClaims)["id"].(float64)

		context.Set(r, "userId", userId)

		// call next to continue
		next(w, r)
	}
}
