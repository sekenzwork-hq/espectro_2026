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

func (v VenueUsecases) CreateVenue(venue entity.VenueCreateEntity) (uuid.UUID, error) {

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

	venueCreateEntity := entity.VenueEntity{
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

func (v VenueUsecases) UpdateVenue(venueId string, newVenue entity.VenueUpdateEntity) error {

	if correct := pkg.ValidateUUID(venueId); !correct {
		return &customerrors.ValidationError{OrgError: "Invalid venue id"}
	}

	if newVenue.Country != nil {

		if correct := pkg.ValidateCountryOrState(*newVenue.Country); !correct {
			return &customerrors.ValidationError{OrgError: "Invalid country name"}
		}
	}

	if newVenue.State != nil {
		if correct := pkg.ValidateCountryOrState(*newVenue.State); !correct {
			return &customerrors.ValidationError{OrgError: "Invalid state name"}
		}
	}

	if newVenue.City != nil {
		if correct := pkg.ValidateCity(*newVenue.City); !correct {
			return &customerrors.ValidationError{OrgError: "Invalid city name"}
		}
	}

	err := v.repo.UpdateVenue(venueId, newVenue.Country, newVenue.State, newVenue.City)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &customerrors.NotFoundError{OrgError: "Venue does not exist"}
	} else if err != nil {
		return &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}
	return nil

}

func (v VenueUsecases) RetrieveVenue(page int, limit int) ([]entity.VenueEntity, error) {

	emptyVenue := []entity.VenueEntity{}
	if limit > 150 {
		return emptyVenue, &customerrors.SizeError{OrgError: "Limit should be less than or equal to 150"}
	}

	offset := pkg.GetOffset(limit, page)

	venue, retrievalErr := v.repo.RetrieveVenue(offset, limit)

	if retrievalErr != nil {
		return emptyVenue, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}

	return venue, nil
}
