package models

type Stats struct {
	Users        string `json:"users"`
	Bookings     string `json:"bookings"`
	OpenBookings string `json:"open_bookings"`
}
