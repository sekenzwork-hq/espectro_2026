package entity

import "time"

type GalleryEntity struct {
	Id        string    `json:"id" gorm:"column:id;type uuid;primaryKey"`
	Name      string    `json:"name" gorm:"column:name"`
	ImageUrls []string  `json:"image_urls" gorm:"column:image_urls"`
	VenueId   string    `json:"venue_id" gorm:"column:venue_id"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;type timestampz;default:now()"`
}
