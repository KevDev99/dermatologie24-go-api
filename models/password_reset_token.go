package models

import "time"

type PasswordResetToken struct {
	Token     string
	UserID    int
	ExpiresAt time.Time
}