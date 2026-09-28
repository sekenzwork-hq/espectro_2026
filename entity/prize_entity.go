package entity

import (
	"espectro/enums"
	"mime/multipart"
	"time"
)

type PrizeDBRetrieveEntity struct {
	Id          string              `json:"id" gorm:"column:id;type uuid"`
	Title       string              `json:"title" gorm:"column:title"`
	Description string              `json:"description" gorm:"column:description"`
	Amount      float32             `json:"amount" gorm:"column:amount"`
	LogoURL     *string             `json:"logo_url" gorm:"column:logo_url"`
	Type        enums.PrizeTypeEnum `json:"type" gorm:"column:type"`
	CreatedAt   time.Time           `json:"created_at" gorm:"column:created_at; type timestampz; default:now()"`
}

type PrizeDBCreateEntity struct {
	Id          string              `gorm:"column:id;type uuid"`
	Title       string              `gorm:"column:title"`
	Description string              `gorm:"column:description"`
	Amount      float32             `gorm:"column:amount"`
	LogoURL     *string             `gorm:"column:logo_url"`
	EventId     string              `gorm:"column:event_id"`
	Type        enums.PrizeTypeEnum `gorm:"column:type"`
}
type PrizeDBUpdateEntity struct {
	Title       *string              `gorm:"column:title"`
	Description *string              `gorm:"column:description"`
	Amount      *float32             `gorm:"column:amount"`
	LogoURL     *string              `gorm:"column:logo_url"`
	EventId     *string              `gorm:"column:event_id"`
	Type        *enums.PrizeTypeEnum `gorm:"column:type"`
}

type PrizeCreateEntity struct {
	Title       string
	Description string
	Amount      float32
	EventId     string
	Logo        *multipart.FileHeader
	Type        enums.PrizeTypeEnum
}

type PrizeUpdateEntity struct {
	Title       *string
	Description *string
	Amount      *float32
	EventId     *string
	Logo        *multipart.FileHeader
	Type        *enums.PrizeTypeEnum
}
