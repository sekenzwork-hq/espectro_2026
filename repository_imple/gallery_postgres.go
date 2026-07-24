package repositoryimple

import (
	"espectro/entity"
	"time"

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

func (g GalleryPostgresRepo) UpdateGallery(newGallery entity.GalleryUpdateEntity) (entity.GalleryEntity, error) {

	var entity entity.GalleryEntity

	out := g.db.Raw(
		`UPDATE gallery SET 
		name=COALESCE(?,name),
		image_urls=COALESCE(?,image_urls)
		
		WHERE id=? AND deleted_at IS NULL
		RETURNING id,name,created_at,image_urls
		`,
		newGallery.Name, newGallery.ImageUrls, newGallery.GallerId,
	).Scan(&entity)

	if out.Error != nil {
		return entity, out.Error
	} else if out.RowsAffected == 0 {
		return entity, gorm.ErrRecordNotFound
	}
	return entity, nil
}

func (g GalleryPostgresRepo) DeleteGallery(galleryId string) error {

	out := g.db.
		Table("gallery").
		Where("id=? AND deleted_at IS NULL", galleryId).
		Update("deleted_at", time.Now().UTC())

	if out.Error != nil {
		return out.Error
	} else if out.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (g GalleryPostgresRepo) RetrieveGalleries(limit int, offset int) ([]entity.GalleryEntity, error) {

	var galleries []entity.GalleryEntity

	err := g.db.
		Table("gallery").
		Select("id,name,created_at,image_urls").
		Where("deleted_at IS NULL").
		Offset(offset).
		Limit(limit).
		Scan(&galleries).Error

	return galleries, err
}

func (g GalleryPostgresRepo) GalleryExists(galleryId string) (bool, error) {

	var exists bool
	err := g.db.Raw(
		`SELECT EXISTS (SELECT 1 FROM gallery WHERE id=? AND deleted_at IS NULL)`,
		galleryId,
	).Scan(&exists).Error

	return exists, err
}
