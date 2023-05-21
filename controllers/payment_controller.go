package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/KevDev99/dermatologie24-go-api/configs"
	"github.com/KevDev99/dermatologie24-go-api/models"
	"github.com/KevDev99/dermatologie24-go-api/utils"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/customer"
	"github.com/stripe/stripe-go/v74/ephemeralkey"
	"github.com/stripe/stripe-go/v74/paymentintent"
)

func StripePaymentSheet() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		stripe.Key = configs.StripeApiKey()

		var paymentRequest models.PaymentRequest

		// Parse the request body
		err := json.NewDecoder(r.Body).Decode(&paymentRequest)

		if err != nil {
			utils.SendResponse(rw, http.StatusInternalServerError, "Internal Error", map[string]interface{}{"data": err.Error()})
		}

		customerParams := &stripe.CustomerParams{Email: &paymentRequest.Email}
		customer, err := customer.New(customerParams)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		ephemeralKeyParams := &stripe.EphemeralKeyParams{
			Customer: stripe.String(customer.ID),

			StripeVersion: stripe.String("2022-11-15"),
		}
		ephemeralKey, err := ephemeralkey.New(ephemeralKeyParams)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		paymentIntentParams := &stripe.PaymentIntentParams{
			Amount:   stripe.Int64(int64(paymentRequest.Amount * 100)),
			Currency: stripe.String(paymentRequest.Currency),
			Customer: stripe.String(customer.ID),
			PaymentMethodTypes: []*string{
				stripe.String("card"),
			},
			/* 	AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
				Enabled: stripe.Bool(true),
			}, */
		}
		paymentIntent, err := paymentintent.New(paymentIntentParams)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		response := struct {
			PaymentIntentId string `json:"paymentIntentId"`
			PaymentIntent   string `json:"paymentIntent"`
			EphemeralKey    string `json:"ephemeralKey"`
			Customer        string `json:"customer"`
			PublishableKey  string `json:"publishableKey"`
		}{
			PaymentIntentId: paymentIntent.ID,
			PaymentIntent:   paymentIntent.ClientSecret,
			EphemeralKey:    ephemeralKey.Secret,
			Customer:        customer.ID,
			PublishableKey:  configs.StripePublishableKey(),
		}

		// Return a success response
		utils.SendResponse(rw, http.StatusOK, "success", map[string]interface{}{"data": response})
	}
}
