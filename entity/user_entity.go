package entity

import "time"

type UserEntity struct {
	Id          string    `json:"id" gorm:"primaryKey"`
	Fullname    string    `json:"fullname" gorm:"column:fullname"`
	Email       string    `json:"email" gorm:"column:email"`
	PhoneNumber int       `json:"phone_number" gorm:"column:phone"`
	Country     string    `json:"country" gorm:"column:country"`
	CountryCode string    `json:"country_code"`
	State       string    `json:"state" gorm:"column:state"`
	City        string    `json:"city" gorm:"column:city"`
	Usertype    string    `json:"user_type" gorm:"user_type"`
	CreatedAt   time.Time `json:"created_at"`
}
