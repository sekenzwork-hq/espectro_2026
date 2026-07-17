package entity

import "time"

type InvestorEntity struct {
	Id          string    `json:"id" gorm:"column:id;type uuid;primaryKey"`
	Name        string    `json:"name" gorm:"column:name"`
	WebsiteUrl  *string   `json:"website_url" gorm:"column:website_url"`
	PhoneNumber string    `json:"phone_number" gorm:"column:phone"`
	Email       string    `json:"email" gorm:"column:email"`
	LogoUrl     *string   `json:"logo_url" gorm:"column:logo_url"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at;type timestampz;default:now()"`
}

type InvestorUpdateEntity struct {
	Name        *string
	WebsiteUrl  *string
	PhoneNumber *string
	Email       *string
	LogoUrl     *string
}
