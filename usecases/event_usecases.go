package usecases

import (
	"errors"
	customerrors "espectro/custom_errors"
	"espectro/entity"
	"espectro/enums"
	"espectro/pkg"
	"espectro/repository"
	"time"

	"gorm.io/gorm"
)

type EventUsecases struct {
	eventRepo    repository.EventRepo
	spectrumRepo repository.SpectrumRepo
	venueRepo    repository.VenueRepo
}

func NewEventUsecases(eventRepo repository.EventRepo, spectrumRepo repository.SpectrumRepo, venueRepo repository.VenueRepo) EventUsecases {
	return EventUsecases{eventRepo: eventRepo, spectrumRepo: spectrumRepo, venueRepo: venueRepo}
}

func (e EventUsecases) CreateEvent(event entity.EventFromJsonEntity) (entity.EventEntity, error) {

	eventToInsert, err := e.validateEventDetailsAndCheckExistence(event.Name,
		event.Description,
		event.SpectrumId,
		event.VenueId,
		event.ParticipantLimit,
		event.StartDate,
		event.EndDate,
		event.EventMode,
		event.EventType,
		event.Status,
		event.ContactEmail,
		event.IsFeatured,
	)

	if err != nil {
		return entity.EventEntity{}, err
	}

	newEvent, insertionErr := e.eventRepo.CreateEvent(eventToInsert)
	if insertionErr != nil {
		return entity.EventEntity{}, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}
	return newEvent, nil

}

func (e EventUsecases) UpdateEvent(event entity.EventUpdateEntity) error {

	if !pkg.ValidateUUID(event.Id) {

		return &customerrors.ValidationError{OrgError: "Invalid event id"}
	}

	_, validationErr := e.validateEventDetailsAndCheckExistence(
		event.Name,
		event.Description,
		event.SpectrumId,
		event.VenueId,
		event.ParticipantLimit,
		event.StartDate,
		event.EndDate,
		event.EventMode,
		event.EventType,
		event.Status,
		event.ContactEmail,
		event.IsFeatured,
	)

	if validationErr != nil {
		return validationErr
	}

	updationErr := e.eventRepo.UpdateEvent(event.Id, event)

	if errors.Is(updationErr, gorm.ErrRecordNotFound) {
		return &customerrors.NotFoundError{OrgError: "Event does not exist"}
	} else if updationErr != nil {
		return &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}

	return nil
}

func (e EventUsecases) validateEventDetailsAndCheckExistence(name string,
	description string,
	spectrumId string,
	venueId string,
	participantLimit *int,
	startDate *string,
	endDate *string,
	eventMode enums.EventMode,
	eventType enums.EventType,
	status enums.EventStatus,
	contactEmail string,
	isFeatured bool,

) (entity.EventEntity, error) {

	emptyEntity := entity.EventEntity{}

	nameErr := pkg.ValidateSpectrumOrEventName(name)
	if nameErr != nil {
		return emptyEntity, nameErr
	}

	descriptionErr := pkg.ValidateSpectrumOrEventDescription(description)
	if descriptionErr != nil {
		return emptyEntity, descriptionErr
	}

	isSpectrumIdCorrect := pkg.ValidateUUID(spectrumId)
	if !isSpectrumIdCorrect {
		return emptyEntity, &customerrors.ValidationError{OrgError: "Invalid spectrum id"}
	}

	if !status.IsValid() {
		return emptyEntity, &customerrors.ValidationError{OrgError: "Invalid event status"}
	}

	if participantLimit != nil && *(participantLimit) < 0 {
		return emptyEntity, &customerrors.ValidationError{OrgError: "Invalid participant limit"}
	}

	var parsedStartTime *time.Time
	if startDate != nil {
		time, err := pkg.ParseTime(*startDate)
		if err != nil {
			return emptyEntity, &customerrors.ValidationError{OrgError: "Invalid start date"}
		}
		parsedStartTime = &time

	}

	var parsedEndTime *time.Time
	if endDate != nil {
		time, err := pkg.ParseTime(*(endDate))
		if err != nil {
			return emptyEntity, &customerrors.ValidationError{OrgError: "Invalid end date"}
		}
		parsedEndTime = &time
	}

	if parsedStartTime != nil && parsedEndTime != nil {

		year, month, day := parsedStartTime.Date()
		now := time.Now()
		if year < now.Year() || month < now.Month() || day < now.Day() {
			return emptyEntity, &customerrors.ValidationError{OrgError: "Start date should be today or after today"}
		} else if parsedStartTime.After(*parsedEndTime) {
			return emptyEntity, &customerrors.ValidationError{OrgError: "Start date should be on same day as end date or before end date"}
		}
	}

	if !eventMode.IsValid() {
		return emptyEntity, &customerrors.ValidationError{OrgError: "Invalid event mode"}
	}

	if !eventType.IsValid() {
		return emptyEntity, &customerrors.ValidationError{OrgError: "Invalid event type"}
	}

	if !pkg.ValidateEmail(contactEmail) {
		return emptyEntity, &customerrors.ValidationError{OrgError: "Invalid contact email"}
	}

	if !pkg.ValidateUUID(venueId) {
		return emptyEntity, &customerrors.ValidationError{OrgError: "Invalid venue id"}
	}

	spectrumExists, spectrumExistenceErr := e.spectrumRepo.CheckSpectrumExists(spectrumId)

	if spectrumExistenceErr != nil {
		return emptyEntity, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	} else if !spectrumExists {
		return emptyEntity, &customerrors.NotFoundError{OrgError: "Spectrum does not exist"}
	}

	venueExists, venueExistanceErr := e.venueRepo.CheckVenueExists(venueId)

	if venueExistanceErr != nil {
		return emptyEntity, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	} else if !venueExists {
		return emptyEntity, &customerrors.NotFoundError{OrgError: "Venue does not exist"}
	}

	eventToInsert := entity.EventEntity{
		Name:             name,
		Description:      description,
		SpectrumId:       spectrumId,
		Status:           status,
		EventMode:        eventMode,
		EventType:        eventType,
		ParticipantLimit: participantLimit,
		StartDate:        parsedStartTime,
		EndDate:          parsedEndTime,
		IsFeatured:       isFeatured,
		ContactEmail:     contactEmail,
		VenueId:          venueId,
	}

	return eventToInsert, nil
}
