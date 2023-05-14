package models

import "reflect"

type User struct {
	Id                  int    `json:"id"`
	Email               string `json:"email" validate:"required"`
	Password            string `json:"password" validate:"required"`
	GotPatientDetailsYN bool   `json:"got_patient_details_yn"`
}

type UserWithoutPassword struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

func (u User) GetJsonFields() []string {
	var jsonFields []string
	val := reflect.ValueOf(u)
	for i := 0; i < val.Type().NumField(); i++ {
		jsonFields = append(jsonFields, val.Type().Field(i).Tag.Get("json"))

	}

	return jsonFields
}
