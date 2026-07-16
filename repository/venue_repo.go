package repository

import (
	"espectro/entity"
)

type VenueRepo interface {
	CreateVenue(venue entity.VenueEntity) (entity.VenueEntity, error)
	DeleteVenue(venueId string) error
	UpdateVenue(venueId string, country *string, state *string, city *string) (entity.VenueEntity, error)
	RetrieveVenue(offset int, limit int) ([]entity.VenueEntity, error)
	CheckVenueExists(venueId string) (bool, error)
}
