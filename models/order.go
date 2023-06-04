package models

import "time"

type Order struct {
	Id            int         `json:"id" gorm:"foreignKey:OrderId;references:Id"`
	UserId        int         `json:"user_id" validate:"required"`
	Message       string      `json:"message" validate:"required"`
	StatusId      int         `json:"status_id" validate:"required"`
	PaymentExtId  string      `json:"payment_ext_id" validate:"required"`
	PaymentTypeId string      `json:"payment_type_id" validate:"required"`
	User          User        `json:"user" gorm:"foreignKey:UserID;references:Id"`
	OrderFile     []OrderFile `json:"order_files"`
}

type OrderFile struct {
	OrderId   int       `json:"order_id" validate:"required"`
	FileId    int       `json:"file_id" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
