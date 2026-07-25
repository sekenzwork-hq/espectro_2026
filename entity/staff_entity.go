package entity

import (
	"espectro/enums"
	"time"
)

type StaffEntity struct {
	Id          string          `json:"id" gorm:"column:id;type uuid;default:gen_random_uuid();primaryKey"`
	Name        string          `json:"fullname" gorm:"column:fullname"`
	Role        enums.StaffRole `json:"role" gorm:"column:role"`
	PhoneNumber string          `json:"phone_number" gorm:"column:phone"`
	Email       string          `json:"email" gorm:"column:email"`
	CreatedAt   time.Time       `json:"created_at" gorm:"column:created_at;type timestampz;default:now()"`
}

type StaffFromJsonEntity struct {
	Name        string          `json:"fullname" `
	Role        enums.StaffRole `json:"role" `
	PhoneNumber string          `json:"phone_number"`
	Email       string          `json:"email"`
}
