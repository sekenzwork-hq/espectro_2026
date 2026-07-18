package entity

import (
	"espectro/enums"
	"time"
)

type SponsorEntity struct {
	Id              string              `json:"id" gorm:"column:id;type uuid;primaryKey"`
	Name            string              `json:"name" gorm:"column:name"`
	Amount          float32             `json:"amount" gorm:"column:amount"`
	ProfileOrOrgUrl *string             `json:"profile_or_org_url" gorm:"column:profile_or_org_url"`
	Sponsored       enums.SponsoredType `json:"sponsored_type" gorm:"column:sponsored_type"`
	CreatedAt       time.Time           `json:"created_at" gorm:"column:created_at;type timestampz;default:now()"`
}

type SponsorUpdateEntity struct {
	Id              string
	Name            *string
	Amount          *float32
	ProfileOrOrgUrl *string
	Sponsored       *enums.SponsoredType
}
