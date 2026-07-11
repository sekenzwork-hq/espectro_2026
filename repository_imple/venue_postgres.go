package repositoryimple

import (
	"espectro/entity"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VenuePostgresRepo struct {
	db *gorm.DB
}

func NewVenuePostgresRepo(db *gorm.DB) VenuePostgresRepo {
	return VenuePostgresRepo{db: db}
}

func (v VenuePostgresRepo) CreateVenue(venue entity.VenueEntity) (uuid.UUID, error) {
	err := v.db.
		Table("venue").
		Create(&venue).Error
	return venue.Id, err
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

func (v VenuePostgresRepo) UpdateVenue(venueId string, country *string, state *string, city *string) error {

	data := map[string]any{}
	if country != nil {
		data["country"] = *country
	}
	if state != nil {
		data["state"] = *state
	}
	if city != nil {
		data["city"] = *city
	}

	out := v.db.
		Table("venue").
		Where("id=? AND deleted_at IS NULL", venueId).
		Updates(data)

	if out.Error != nil {
		return out.Error
	} else if out.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
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
