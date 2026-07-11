package repository

import (
	"espectro/entity"

	"github.com/google/uuid"
)

type VenueRepo interface {
	CreateVenue(venue entity.VenueEntity) (uuid.UUID, error)
	DeleteVenue(venueId string) error
	UpdateVenue(venueId string, country *string, state *string, city *string) error
	RetrieveVenue(offset int, limit int) ([]entity.VenueEntity, error)
}
