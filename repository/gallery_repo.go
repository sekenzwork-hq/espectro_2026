package repository

import "espectro/entity"

type GalleryRepo interface {
	CreateGallery(gallery entity.GalleryEntity) (entity.GalleryEntity, error)
	UpdateGallery(newGallery entity.GalleryUpdateEntity) (entity.GalleryEntity, error)
	DeleteGallery(galleryId string) error
	RetrieveGalleries(limit int, offset int) ([]entity.GalleryEntity, error)
}
