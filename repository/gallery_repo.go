package repository

import "espectro/entity"

type GalleryRepo interface {
	CreateGallery(gallery entity.GalleryEntity) (entity.GalleryEntity, error)
}
