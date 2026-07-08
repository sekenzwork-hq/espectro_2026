package entity

import "github.com/google/uuid"

type AdminEnteredLoginCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AdminDBLoginCredentials struct {
	Id       uuid.UUID `gorm:"column:id"`
	Email    string    `gorm:"column:email"`
	Password string    `gorm:"column:password"`
}

type AdminCreateEntity struct {
	Id       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Email    string    `gorm:"column:email"`
	Fullname string    `gorm:"column:fullname"`
	Password string    `gorm:"column:password"`
	Role     string    `gorm:"column:admin_role"`
}
