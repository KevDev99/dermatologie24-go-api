package models

import (
	"reflect"
	"time"
)

type ProfileDetails struct {
	UserId      int       `json:"user_id" gorm:"primaryKey"`
	FullName    string    `json:"full_name" validate:"required"`
	BirthDate   string    `json:"birth_date" validate:"required"`
	PhoneNumber string    `json:"phone_number" validate:"required" `
	Gender      string    `json:"gender" validate:"required"`
	Street      string    `json:"street" validate:"required"`
	ZipCode     int       `json:"zip_code" validate:"required"`
	Location    string    `json:"location" validate:"required"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (u ProfileDetails) GetJsonFields() []string {
	var jsonFields []string
	val := reflect.ValueOf(u)
	for i := 0; i < val.Type().NumField(); i++ {
		jsonFields = append(jsonFields, val.Type().Field(i).Tag.Get("json"))

	}

	return jsonFields
}
