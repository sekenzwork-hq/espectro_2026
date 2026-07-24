package repository

import "espectro/entity"

type EventGalleryRepo interface {
	AddGalleryToEvent(eventAndGallery entity.EventGalleryEntity) error
	IsGalleryAddedToEvent(eventAndGallery entity.EventGalleryEntity) (bool, error)
	DeleteGalleryFromAnEvent(eventAndGallery entity.EventGalleryEntity) error
	DeleteEventOrGallery(eventOrGalleryId string) error
}
