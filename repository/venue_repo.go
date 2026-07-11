package repository

import (
	"espectro/entity"

	"github.com/google/uuid"
)

type VenueRepo interface {
	CreateVenue(venue entity.VenueCreateEntity) (uuid.UUID, error)
}
