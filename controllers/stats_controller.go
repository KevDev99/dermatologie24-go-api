package controllers

import (
	"net/http"

	"github.com/KevDev99/dermatologie24-go-api/configs"
	"github.com/KevDev99/dermatologie24-go-api/utils"
)

func Stats() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		type Result struct {
			Field string
			Total int
		}

		var results []Result

		configs.DB.Raw("SELECT 'orders' AS field, COUNT(*) AS total FROM orders " +
			"UNION " +
			"SELECT 'users' AS field, COUNT(*) AS total FROM users " +
			"UNION " +
			"SELECT 'open_orders' AS field, COUNT(*) AS total FROM orders WHERE status_id=0").
			Scan(&results)

		// Return a success response
		utils.SendResponse(rw, http.StatusOK, "success", map[string]interface{}{"stats": results})
	}
}
