package repositoryimple

import (
	"espectro/entity"
	"time"

	"gorm.io/gorm"
)

type VenuePostgresRepo struct {
	db *gorm.DB
}

func NewVenuePostgresRepo(db *gorm.DB) VenuePostgresRepo {
	return VenuePostgresRepo{db: db}
}

func (v VenuePostgresRepo) CreateVenue(venue entity.VenueEntity) (entity.VenueEntity, error) {
	err := v.db.
		Table("venue").
		Create(&venue).Error
	return venue, err
}

func (v VenuePostgresRepo) DeleteVenue(venueId string) error {

	out := v.db.
		Table("venue").
		Where("id=? AND deleted_at IS NULL", venueId).
		Update("deleted_at", time.Now().UTC())

	if out.Error != nil {
		return out.Error
	} else if out.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (v VenuePostgresRepo) UpdateVenue(venueId string, country *string, state *string, city *string) (entity.VenueEntity, error) {

	var entity entity.VenueEntity
	out := v.db.
		Raw(`UPDATE venue SET 
	country=COALESCE(?,country),
	state=COALESCE(?,state),
	city=COALESCE(?,city) 
	WHERE id=? AND deleted_at IS NULL
	RETURNING id,country,state,city,created_at
	`, country, state, city, venueId,
		).Scan(&entity)

	if out.Error != nil {
		return entity, out.Error
	} else if out.RowsAffected == 0 {
		return entity, gorm.ErrRecordNotFound
	}

	return entity, nil
}

func (v VenuePostgresRepo) RetrieveVenue(offset int, limit int) ([]entity.VenueEntity, error) {

	var venue []entity.VenueEntity
	out := v.db.
		Table("venue").
		Select("id,country,state,city,created_at").
		Where("deleted_at IS NULL").
		Offset(offset).
		Limit(limit).
		Scan(&venue)

	return venue, out.Error
}

func (v VenuePostgresRepo) CheckVenueExists(venueId string) (bool, error) {

	var exists bool
	err := v.db.
		Raw("SELECT EXISTS (SELECT 1 FROM venue WHERE id=? AND deleted_at IS NULL)", venueId).
		Scan(&exists).Error
	return exists, err
}
