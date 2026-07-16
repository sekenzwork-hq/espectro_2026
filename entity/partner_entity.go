package entity

import "time"

type PartnerEntity struct {
	Id        string    `json:"id" gorm:"primaryKey;type uuid"`
	Name      string    `json:"name" gorm:"column:name"`
	LogoUrl   *string   `json:"logo_url" gorm:"column:logo_url"`
	CreatedAt time.Time `json:"created_at" gorm:"type timestampz;default:now()"`
}

type PartnerCreateEntity struct {
	Name    string  `json:"name"`
	LogoUrl *string `json:"logo_url"`
}

type PartnerUpdateEntity struct {
	Id      string  `json:"id"`
	Name    *string `json:"name" gorm:"column:name"`
	LogoUrl *string `json:"logo_url" gorm:"column:logo_url"`
}
