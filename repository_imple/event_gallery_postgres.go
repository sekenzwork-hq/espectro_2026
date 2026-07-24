package repositoryimple

import (
	"espectro/entity"
	"time"

	"gorm.io/gorm"
)

type EventGalleryPostgresRepo struct {
	db *gorm.DB
}

func NewEventGalleryPostgresRepo(db *gorm.DB) EventGalleryPostgresRepo {
	return EventGalleryPostgresRepo{db: db}
}

func (e EventGalleryPostgresRepo) AddGalleryToEvent(eventAndGallery entity.EventGalleryEntity) error {
	err := e.db.
		Table("events_gallery").
		Create(&eventAndGallery).Error

	return err
}

func (e EventGalleryPostgresRepo) IsGalleryAddedToEvent(eventAndGallery entity.EventGalleryEntity) (bool, error) {
	var exists bool

	err := e.db.Raw(
		`SELECT EXISTS (SELECT 1 FROM events_gallery WHERE event_id=? AND gallery_id=? AND deleted_at IS NULL)`,
		eventAndGallery.EventId, eventAndGallery.GalleryId,
	).Scan(&exists).Error

	return exists, err
}
func (e EventGalleryPostgresRepo) DeleteGalleryFromAnEvent(eventAndGallery entity.EventGalleryEntity) error {

	out := e.db.
		Table("events_gallery").
		Where("event_id=? AND gallery_id=? AND deleted_at IS NULL", eventAndGallery.EventId, eventAndGallery.GalleryId).
		Update("deleted_at", time.Now().UTC())

	if out.Error != nil {
		return out.Error
	} else if out.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (e EventGalleryPostgresRepo) DeleteEventOrGallery(eventOrGalleryId string) error {

	out := e.db.
		Table("events_gallery").
		Where("event_id=? OR gallery_id=? AND deleted_at IS NULL", eventOrGalleryId, eventOrGalleryId).
		Update("deleted_at", time.Now().UTC())

	if out.Error != nil {
		return out.Error
	} else if out.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
