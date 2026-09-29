package usecases

import (
	"errors"
	customerrors "espectro/custom_errors"
	"espectro/entity"
	"espectro/pkg"
	"espectro/repository"

	"gorm.io/gorm"
)

type VenueUsecases struct {
	repo repository.VenueRepo
}

func NewVenueUsecases(repo repository.VenueRepo) VenueUsecases {
	return VenueUsecases{repo: repo}
}

func (v VenueUsecases) CreateVenue(venue entity.VenueCreateEntity) (entity.VenueEntity, error) {

	emptyVenue := entity.VenueEntity{}
	isCountryCorrect := pkg.ValidateCountryOrState(venue.Country)

	if !isCountryCorrect {
		return emptyVenue, &customerrors.ValidationError{DisplayError: "Invalid country name"}
	}

	isStateCorrect := pkg.ValidateCountryOrState(venue.State)

	if !isStateCorrect {
		return emptyVenue, &customerrors.ValidationError{DisplayError: "Invalid state name"}
	}

	isCityCorrect := pkg.ValidateCity(venue.City)

	if !isCityCorrect {
		return emptyVenue, &customerrors.ValidationError{DisplayError: "Invalid city name"}
	}

	venueCreateEntity := entity.VenueEntity{
		Country: venue.Country,
		State:   venue.State,
		City:    venue.City,
	}

	insertedVenue, venueErr := v.repo.CreateVenue(venueCreateEntity)

	if venueErr != nil {
		return emptyVenue, &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
	}

	return insertedVenue, nil

}

func (v VenueUsecases) DeleteVenue(venueId string) error {

	if correct := pkg.ValidateUUID(venueId); !correct {
		return &customerrors.ValidationError{DisplayError: "Invalid venue id"}
	}

	deletionErr := v.repo.DeleteVenue(venueId)

	if errors.Is(deletionErr, gorm.ErrRecordNotFound) {
		return &customerrors.NotFoundError{DisplayError: "Venue does not exist"}
	} else if deletionErr != nil {
		return &customerrors.ServerError{DisplayError: "Something went wrong"}
	}

	return nil
}

func (v VenueUsecases) UpdateVenue(venueId string, newVenue entity.VenueUpdateEntity) (entity.VenueEntity, error) {

	emptyVenue := entity.VenueEntity{}
	if correct := pkg.ValidateUUID(venueId); !correct {
		return emptyVenue, &customerrors.ValidationError{DisplayError: "Invalid venue id"}
	}

	if newVenue.Country != nil {

		if correct := pkg.ValidateCountryOrState(*newVenue.Country); !correct {
			return emptyVenue, &customerrors.ValidationError{DisplayError: "Invalid country name"}
		}
	}

	if newVenue.State != nil {
		if correct := pkg.ValidateCountryOrState(*newVenue.State); !correct {
			return emptyVenue, &customerrors.ValidationError{DisplayError: "Invalid state name"}
		}
	}

	if newVenue.City != nil {
		if correct := pkg.ValidateCity(*newVenue.City); !correct {
			return emptyVenue, &customerrors.ValidationError{DisplayError: "Invalid city name"}
		}
	}

	updateVenue, err := v.repo.UpdateVenue(venueId, newVenue.Country, newVenue.State, newVenue.City)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return updateVenue, &customerrors.NotFoundError{DisplayError: "Venue does not exist"}
	} else if err != nil {
		return updateVenue, &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
	}
	return updateVenue, nil

}

func (v VenueUsecases) RetrieveVenue(page int) ([]entity.VenueEntity, error) {

	emptyVenue := []entity.VenueEntity{}

	offset := pkg.GetOffset(50, page)

	venue, retrievalErr := v.repo.RetrieveVenue(offset)

	if retrievalErr != nil {
		return emptyVenue, &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
	}

	return venue, nil
}
