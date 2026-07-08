package entity

type AdminEnteredLoginCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AdminDBLoginCredentials struct {
	Email    string `gorm:"column:email"`
	Password string `gorm:"column:password"`
}
