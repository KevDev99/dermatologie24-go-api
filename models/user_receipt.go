package models

import (
	"reflect"
	"time"
)

type UserRecipe struct {
	Id          int       `json:"id"`
	UserId      string    `json:"user_id" validate:"required"`
	Name        string    `json:"name" validate:"required"`
	Information string    `json:"information" validate:"required"`
	FilePath    string    `json:"file_path"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (u UserRecipe) GetJsonFields() []string {
	var jsonFields []string
	val := reflect.ValueOf(u)
	for i := 0; i < val.Type().NumField(); i++ {
		jsonFields = append(jsonFields, val.Type().Field(i).Tag.Get("json"))
	}

	return jsonFields
}
