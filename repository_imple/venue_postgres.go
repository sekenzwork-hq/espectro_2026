package repositoryimple

import (
	"espectro/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VenuePostgresRepo struct {
	db *gorm.DB
}

func NewVenuePostgresRepo(db *gorm.DB) VenuePostgresRepo {
	return VenuePostgresRepo{db: db}
}

func (v VenuePostgresRepo) CreateVenue(venue entity.VenueCreateEntity) (uuid.UUID, error) {
	err := v.db.
		Table("venue").
		Create(&venue).Error
	return venue.Id, err
}
