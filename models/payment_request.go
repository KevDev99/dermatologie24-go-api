package models

type PaymentRequest struct {
	Amount   int    `json:"amount"`
	Currency string `json:"currency"`
	Email    string `json:"email"`
}
