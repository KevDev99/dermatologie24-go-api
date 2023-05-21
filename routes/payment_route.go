package routes

import (
	"github.com/KevDev99/dermatologie24-go-api/controllers"
	"github.com/KevDev99/dermatologie24-go-api/middleware"
	"github.com/gorilla/mux"
)

func PaymentRoute(router *mux.Router) {
	router.Handle("/payment/init-paymentsheet", middleware.AuthMiddleware(controllers.StripePaymentSheet())).Methods("POST")
}
