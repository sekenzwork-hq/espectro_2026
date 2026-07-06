package domain

import "time"

type Usertype int

const (
	Entrepreneur = iota
	Student
	Employee
	Other
)

type User struct {
	Id        string
	Fullname  string
	Email     string
	Country   string
	State     string
	City      string
	Usertype  Usertype
	CreatedAt time.Time
}

type UserRepository interface {
	RegisterUser() (string, error)
	GetUserById(string) (User, error)
	GetAllUsers()
}
