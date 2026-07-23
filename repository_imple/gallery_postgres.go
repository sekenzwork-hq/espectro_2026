package repositoryimple

import (
	"espectro/entity"

	"gorm.io/gorm"
)

type GalleryPostgresRepo struct {
	db *gorm.DB
}

func NewGalleryPostgresRepo(db *gorm.DB) GalleryPostgresRepo {
	return GalleryPostgresRepo{db: db}
}

func (g GalleryPostgresRepo) CreateGallery(gallery entity.GalleryEntity) (entity.GalleryEntity, error) {

	err := g.db.
		Table("gallery").
		Create(&gallery).Error

	return gallery, err
}
