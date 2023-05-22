package models

import (
	"reflect"
	"time"
)

type User struct {
	Id                  int         `json:"id"`
	Email               string      `json:"email" validate:"required"`
	Password            string      `json:"password" validate:"required"`
	GotPatientDetailsYN bool        `json:"got_patient_details_yn"`
	EmailConfirmedYN    bool        `json:"email_confirmed_yn"`
	UserDetails         UserDetails `json:"user_details" gorm:"foreignKey:UserID;references:Id"`
}

type UserDetails struct {
	UserID      int       `json:"user_id" gorm:"primaryKey"`
	Firstname   string    `json:"firstname"`
	Lastname    string    `json:"lastname"`
	Street      string    `json:"street"`
	ZipCode     string    `json:"zip_code"`
	CountryCode string    `json:"country_code"`
	PhoneNumber string    `json:"phone_number"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type UserWithoutPassword struct {
	Id                  int    `json:"id"`
	Email               string `json:"email"`
	GotPatientDetailsYN bool   `json:"got_patient_details_yn"`
	EmailConfirmedYN    bool   `json:"email_confirmed_yn"`
}

func (u User) GetJsonFields() []string {
	var jsonFields []string
	val := reflect.ValueOf(u)
	for i := 0; i < val.Type().NumField(); i++ {
		jsonFields = append(jsonFields, val.Type().Field(i).Tag.Get("json"))

	}

	return jsonFields
}
