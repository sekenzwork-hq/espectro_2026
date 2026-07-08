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
	Email    string    `json:"email" gorm:"column:email"`
	Fullname string    `json:"fullname" gorm:"column:fullname"`
	Password string    `json:"password" gorm:"column:password"`
	Role     string    `json:"admin_role" gorm:"column:admin_role"`
}
