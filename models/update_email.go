package models

import "time"

type UpdateEmail struct {
	Email       string
	UserID      int
	ConfirmedYN bool
	ExpiresAt   time.Time
}
