package entity

type EventGalleryEntity struct {
	EventId   string `json:"event_id" gorm:"column:event_id"`
	GalleryId string `json:"gallery_id" gorm:"column:gallery_id"`
}
