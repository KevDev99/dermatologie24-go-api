package models

import "time"

type Recipe struct {
	Id          int       `json:"id" gorm:"foreignKey:OrderId;references:Id"`
	Title       string    `json:"title" validate:"required"`
	Information string    `json:"information" validate:"required"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
