package models

type Order struct {
	Id            int         `json:"id" gorm:"foreignKey:OrderId;references:Id"`
	UserId        int         `json:"user_id" validate:"required"`
	Message       string      `json:"message" validate:"required"`
	StatusId      int         `json:"status_id" validate:"required"`
	PaymentExtId  string      `json:"payment_ext_id" validate:"required"`
	PaymentTypeId string      `json:"payment_type_id" validate:"required"`
	BookingFiles  []OrderFile `json:"booking_files"`
}

type OrderFile struct {
	Id       int `json:"id"`
	OrderId  int `json:"order_id" validate:"required"`
	Name     string
	FilePath string
}
