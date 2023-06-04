package models

import "time"

type File struct {
	Id        int       `json:"id" validate:"required"`
	Filename  string    `json:"filename" validate:"required"`
	FilePath  string    `json:"file_path" validate:"required"`
	FileSize  int       `json:"file_size" validate:"required"`
	User      User      `json:"user" gorm:"foreignKey:UserID;references:Id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
