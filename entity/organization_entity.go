package entity

import (
	"espectro/enums"
	"mime/multipart"
	"time"
)

type OrganizationEntity struct {
	Id           string                   `json:"id" gorm:"column:id;type uuid"`
	Email        string                   `json:"email" gorm:"column:email"`
	LogoUrl      *string                  `json:"logo_url" gorm:"column:logo_url"`
	Status       enums.OrganizationStatus `json:"status" gorm:"column:status"`
	ApprovedBy   *string                  `json:"approved_by" gorm:"approved_by"`
	Fullname     string                   `json:"fullname" gorm:"column:fullname"`
	PhoneNumber  string                   `json:"phone_number" gorm:"column:phone"`
	WebsiteUrl   *string                  `json:"website_url" gorm:"column:website_url"`
	Industry     string                   `json:"industry" gorm:"column:industry"`
	HeadQuarters string                   `json:"headquarters" gorm:"column:headquarters"`
	CreatedAt    time.Time                `json:"created_at" gorm:"column:created_at;type timestampz; default:now()"`
}

type OrganizationCreateEntity struct {
	Id           string
	Email        string
	Logo         *multipart.FileHeader
	Status       enums.OrganizationStatus
	ApprovedBy   *string
	Fullname     string
	PhoneNumber  string
	WebsiteUrl   *string
	Industry     string
	HeadQuarters string
	CreatedAt    time.Time
}
