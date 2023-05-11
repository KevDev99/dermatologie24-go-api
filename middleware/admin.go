package middleware

import (
	"net/http"

	"github.com/KevDev99/dermatologie24-go-api/configs"
	"github.com/gorilla/context"
)

func AdminMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userId := context.Get(r, "userId")

		// check if user id is given
		if userId == nil || userId == "" {
			http.Error(w, "user id not provided", http.StatusUnauthorized)
			return
		}

		type Result struct {
			AdminYN int `json:"adminYN"`
		}

		var result Result
		configs.DB.Raw("SELECT adminYN as admin_yn FROM users WHERE id = ?", userId).Scan(&result)

		if result.AdminYN != 1 {
			http.Error(w, "user not admin", http.StatusUnauthorized)
			return
		}

		// call next to continue
		next(w, r)
	}
}
