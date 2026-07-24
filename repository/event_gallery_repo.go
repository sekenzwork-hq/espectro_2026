package repository

import "espectro/entity"

type EventGalleryRepo interface {
	AddGalleryToEvent(eventAndGallery entity.EventGalleryEntity) error
}
