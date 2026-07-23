package entity

type GalleryEntity struct {
	Id        string   `json:"id" gorm:"column:id;type uuid;primaryKey"`
	Name      string   `json:"name" gorm:"column:name"`
	ImageUrls []string `json:"image_urls" gorm:"column:image_urls"`
	EventId   string   `json:"event_id" gorm:"column:event_id"`
	VenueId   string   `json:"venue_id" gorm:"column:venue_id"`
	CreatedAt string   `json:"created_at" gorm:"column:created_at"`
}
