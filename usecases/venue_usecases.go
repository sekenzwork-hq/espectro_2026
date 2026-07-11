package usecases

import (
	"errors"
	customerrors "espectro/custom_errors"
	"espectro/entity"
	"espectro/pkg"
	"espectro/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VenueUsecases struct {
	repo repository.VenueRepo
}

func NewVenueUsecases(repo repository.VenueRepo) VenueUsecases {
	return VenueUsecases{repo: repo}
}

func (v VenueUsecases) CreateVenue(venue entity.VenueFromJson) (uuid.UUID, error) {

	emptyUUID := uuid.UUID{}
	isCountryCorrect := pkg.ValidateCountryOrState(venue.Country)

	if !isCountryCorrect {
		return emptyUUID, &customerrors.ValidationError{OrgError: "Invalid country name"}
	}

	isStateCorrect := pkg.ValidateCountryOrState(venue.State)

	if !isStateCorrect {
		return emptyUUID, &customerrors.ValidationError{OrgError: "Invalid state name"}
	}

	isCityCorrect := pkg.ValidateCity(venue.City)

	if !isCityCorrect {
		return emptyUUID, &customerrors.ValidationError{OrgError: "Invalid city name"}
	}

	venueCreateEntity := entity.VenueCreateEntity{
		Country: venue.Country,
		State:   venue.State,
		City:    venue.City,
	}

	venueId, venueErr := v.repo.CreateVenue(venueCreateEntity)

	if venueErr != nil {
		return emptyUUID, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}

	return venueId, nil

}

func (v VenueUsecases) DeleteVenue(venueId string) error {

	if correct := pkg.ValidateUUID(venueId); !correct {
		return &customerrors.ValidationError{OrgError: "Invalid venue id"}
	}

	deletionErr := v.repo.DeleteVenue(venueId)

	if errors.Is(deletionErr, gorm.ErrRecordNotFound) {
		return &customerrors.NotFoundError{OrgError: "Venue does not exist"}
	} else if deletionErr != nil {
		return &customerrors.ServerError{OrgError: "Something went wrong"}
	}

	return nil
}
