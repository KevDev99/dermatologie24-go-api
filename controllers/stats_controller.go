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

		configs.DB.Raw("SELECT 'bookings' AS field, COUNT(*) AS total FROM bookings " +
			"UNION " +
			"SELECT 'users' AS field, COUNT(*) AS total FROM users " +
			"UNION " +
			"SELECT 'open_bookings' AS field, COUNT(*) AS total FROM bookings WHERE booking_status=0").
			Scan(&results)

		// Return a success response
		utils.SendResponse(rw, http.StatusOK, "success", map[string]interface{}{"stats": results})
	}
}
