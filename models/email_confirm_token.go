package models

import "time"

type EmailConfirmToken struct {
	Token     string
	UserID    int
	ExpiresAt time.Time
}