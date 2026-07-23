package entity

import (
	"time"

	"github.com/lib/pq"
)

type GalleryEntity struct {
	Id        string         `json:"id" gorm:"column:id;type uuid;primaryKey"`
	Name      string         `json:"name" gorm:"column:name"`
	ImageUrls pq.StringArray `json:"image_urls" gorm:"column:image_urls"`
	CreatedAt time.Time      `json:"created_at" gorm:"column:created_at;type timestampz;default:now()"`
}

type GalleryUpdateEntity struct {
	GallerId  string
	Name      *string
	ImageUrls *pq.StringArray
	VenueId   *string
}
