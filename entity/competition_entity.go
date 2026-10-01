package entity

import (
	"espectro/enums"
	"mime/multipart"
	"time"
)

type CompetitionDBRetrieveEntity struct {
	Id              string                      `json:"id" gorm:"column:id; type uuid"`
	Title           string                      `json:"title" gorm:"column:title"`
	Description     string                      `json:"description" gorm:"column:description"`
	Rules           string                      `json:"rules" gorm:"column:rules"`
	Status          enums.CompetitionStatusEnum `json:"status" gorm:"column:status"`
	OpeningTime     time.Time                   `json:"opening_time" gorm:"column:opening_time"`
	ClosingTime     time.Time                   `json:"closing_time" gorm:"column:closing_time"`
	ResultTime      time.Time                   `json:"result_time" gorm:"column:result_time"`
	RegistrationFee float32                     `json:"registration_fee" gorm:"column:registration_fee"`
	Prize           PrizeDBRetrieveEntity       `json:"prize"`
	Venue           VenueEntity                 `json:"venue"`
	LogoURL         *string                     `json:"logo_url" gorm:"column:logo_url"`
	CreatedAt       time.Time                   `json:"created_at" gorm:"column:created_at; type:timestampz; default:now()"`
}

type CompetitionDBCreateEntity struct {
	Id              string                      `gorm:"column:id; type uuid"`
	EventId         string                      `gorm:"column:event_id"`
	Title           string                      `gorm:"column:title"`
	Description     string                      `gorm:"column:description"`
	Rules           string                      `gorm:"column:rules"`
	Status          enums.CompetitionStatusEnum `gorm:"column:status"`
	OpeningTime     time.Time                   `gorm:"column:opening_time"`
	ClosingTime     time.Time                   `gorm:"column:closing_time"`
	ResultTime      time.Time                   `gorm:"column:result_time"`
	RegistrationFee float32                     `gorm:"column:reg_fee"`
	PrizeId         string                      `gorm:"column:prize_id"`
	VenueId         string                      `gorm:"column:venue_id"`
	LogoURL         *string                     `gorm:"column:logo_url"`
}

type CompetitionDBUpdateEntity struct {
	EventId         *string                      `gorm:"column:event_id"`
	Title           *string                      `gorm:"column:title"`
	Description     *string                      `gorm:"column:description"`
	Rules           *string                      `gorm:"column:rules"`
	Status          *enums.CompetitionStatusEnum `gorm:"column:status"`
	OpeningTime     *time.Time                   `gorm:"column:opening_time"`
	ClosingTime     *time.Time                   `gorm:"column:closing_time"`
	ResultTime      *time.Time                   `gorm:"column:result_time"`
	RegistrationFee *float32                     `gorm:"column:reg_fee"`
	PrizeId         *string                      `gorm:"column:prize_id"`
	VenueId         *string                      `gorm:"column:venue_id"`
	LogoURL         *string                      `gorm:"column:logo_url"`
}

type CompetitionCreateEntity struct {
	EventId         string
	Title           string
	Description     string
	Rules           string
	Status          enums.CompetitionStatusEnum
	OpeningTime     time.Time
	ClosingTime     time.Time
	ResultTime      time.Time
	RegistrationFee float32
	PrizeId         string
	VenueId         string
	Logo            *multipart.FileHeader
}

type CompetitionUpdateEntity struct {
	EventId         *string
	Title           *string
	Description     *string
	Rules           *string
	Status          *enums.CompetitionStatusEnum
	OpeningTime     *time.Time
	ClosingTime     *time.Time
	ResultTime      *time.Time
	RegistrationFee *float32
	PrizeId         *string
	VenueId         *string
	Logo            *multipart.FileHeader
}
