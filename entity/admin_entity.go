package entity

import (
	"time"

	"github.com/google/uuid"
)

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

type AdminEntity struct {
	Id        uuid.UUID `json:"id" gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey"`
	Email     string    `json:"email" gorm:"column:email"`
	Fullname  string    `json:"fullname" gorm:"column:fullname"`
	Role      string    `json:"admin_role" gorm:"column:admin_role"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;type timestampz;default:now()"`
}

type AdminUpdateEntity struct {
	Email    *string `json:"email"`
	Fullname *string `json:"fullname"`
}

type AdminRoleUpdateEntity struct {
	AdminId string `json:"admin_id"`
	Role    string `json:"admin_role"`
}
