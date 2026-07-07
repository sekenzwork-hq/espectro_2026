package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserEntity struct {
	Id          uuid.UUID `json:"id"  gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Fullname    string    `json:"fullname" gorm:"column:fullname"`
	Email       string    `json:"email" gorm:"column:email"`
	PhoneNumber string    `json:"phone_number" gorm:"column:phone"`
	Country     string    `json:"country" gorm:"column:country"`
	CountryCode string    `json:"country_code" gorm:"-"`
	State       string    `json:"state" gorm:"column:state"`
	City        string    `json:"city" gorm:"column:city"`
	Usertype    string    `json:"user_type" gorm:"column:user_type"`
	CreatedAt   time.Time `json:"created_at" gorm:"-"`
}
