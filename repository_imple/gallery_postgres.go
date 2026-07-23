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
		venue_id=COALESCE(?,venue_id),
		image_urls=COALESCE(?,image_urls)
		
		WHERE id=? AND deleted_at IS NULL
		RETURNING id,name,venue_id,created_at,image_urls
		`,
		newGallery.Name, newGallery.VenueId, newGallery.ImageUrls, newGallery.GallerId,
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

func (g GalleryPostgresRepo) RetrieveGalleries(limit int, offset int) ([]entity.GalleryWithVenueEntity, error) {

	var galleries []entity.GalleryWithVenueEntity

	err := g.db.
		Table("gallery g").
		Select("g.id,g.name,g.venue_id,g.created_at,g.image_urls, v.country as country, v.state as state, v.city as city").
		Joins("JOIN venue v ON v.id = g.venue_id AND v.deleted_at IS NULL").
		Where("g.deleted_at IS NULL").
		Offset(offset).
		Limit(limit).
		Scan(&galleries).Error

	return galleries, err
}
