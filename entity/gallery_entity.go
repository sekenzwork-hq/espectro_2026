package entity

import (
	"time"

	"github.com/lib/pq"
)

type GalleryEntity struct {
	Id        string         `json:"id" gorm:"column:id;type uuid;primaryKey"`
	Name      string         `json:"name" gorm:"column:name"`
	ImageUrls pq.StringArray `json:"image_urls" gorm:"column:image_urls"`
	VenueId   string         `json:"venue_id" gorm:"column:venue_id"`
	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at;type timestampz;default:now()"`
}

type GalleryUpdateEntity struct {
	GallerId  string
	Name      *string
	ImageUrls *pq.StringArray
	VenueId   *string
}

type GalleryWithVenueEntity struct {
	Id           string         `json:"id" gorm:"column:id;type uuid;primaryKey"`
	Name         string         `json:"name" gorm:"column:name"`
	ImageUrls    pq.StringArray `json:"image_urls" gorm:"column:image_urls"`
	VenueId      string         `json:"venue_id" gorm:"column:venue_id"`
	VenueCountry *string        `json:"country" gorm:"column:country"`
	VenueState   *string        `json:"state" gorm:"column:state"`
	VenueCity    *string        `json:"city" gorm:"column:city"`
	CreatedAt    time.Time      `json:"created_at" gorm:"column:created_at;type timestampz;default:now()"`
}
