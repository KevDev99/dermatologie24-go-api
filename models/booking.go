package models

type Booking struct {
	Id           int           `json:"id" gorm:"foreignKey:BookingId;references:Id"`
	UserId       int           `json:"user_id" validate:"required"`
	Subject      string        `json:"subject" validate:"required"`
	Message      string        `json:"message" validate:"required"`
	BookingFiles []BookingFile `json:"booking_files"`
}

type BookingFile struct {
	Id        int `json:"id"`
	BookingId int `json:"booking_id" validate:"required"`
	Name      string
	FilePath  string
}
