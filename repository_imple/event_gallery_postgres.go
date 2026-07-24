package repositoryimple

import (
	"espectro/entity"

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
