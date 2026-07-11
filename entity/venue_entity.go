package entity

import (
	"time"

	"github.com/google/uuid"
)

type VenueCreateEntity struct {
	Id        uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Country   string    `json:"country" gorm:"column:country"`
	State     string    `json:"state" gorm:"column:state"`
	City      string    `json:"city" gorm:"column:city"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
}

type VenueFromJson struct {
	Country string `json:"country"`
	State   string `json:"state"`
	City    string `json:"city"`
}
