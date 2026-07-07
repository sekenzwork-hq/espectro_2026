package entity

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
