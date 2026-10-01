package entity

import (
	"time"
)

type UserEntity struct {
	Id          string    `json:"id"  gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Fullname    string    `json:"fullname" gorm:"column:fullname"`
	Username    string    `json:"username" gorm:"column:username"`
	Email       string    `json:"email" gorm:"column:email"`
	PhoneNumber string    `json:"phone_number" gorm:"column:phone"`
	Country     string    `json:"country" gorm:"column:country"`
	CountryCode string    `json:"country_code" gorm:"column:country_code"`
	State       string    `json:"state" gorm:"column:state"`
	City        string    `json:"city" gorm:"column:city"`
	Usertype    string    `json:"user_type" gorm:"column:user_type"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at;type:timestampz;default:now()"`
}

type UserDBCreateEntity struct {
	Username    string `gorm:"column:username"`
	Password    string `gorm:"column:password"`
	Fullname    string `gorm:"column:fullname"`
	Email       string `gorm:"column:email"`
	PhoneNumber string `gorm:"column:phone"`
	Country     string `gorm:"column:country"`
	CountryCode string `json:"country_code" gorm:"column:country_code"`
	State       string `gorm:"column:state"`
	City        string `gorm:"column:city"`
	Usertype    string `gorm:"column:user_type"`
}

type UserDBCredentialsEntity struct {
	Id       string `gorm:"colum:id"`
	Username string `gorm:"column:username"`
	Password string `gorm:"column:password"`
}

type UserCredentialsJsonEntity struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type UserJsonCreateEntity struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	Fullname    string `json:"fullname"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone"`
	Country     string `json:"country"`
	CountryCode string `json:"country_code"`
	State       string `json:"state"`
	City        string `json:"city"`
	Usertype    string `json:"user_type"`
}
