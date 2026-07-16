package usecases

import (
	"errors"
	customerrors "espectro/custom_errors"
	"espectro/entity"
	"espectro/enums"
	"espectro/pkg"
	"espectro/repository"
	"fmt"
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

func (e EventUsecases) CreateEvent(event entity.EventCreateEntity) (entity.EventEntity, error) {

	err := e.validateEventDetailsAndCheckExistence(&event.Name,
		&event.Description,
		&event.SpectrumId,
		&event.VenueId,
		event.ParticipantLimit,
		event.StartDate,
		event.EndDate,
		&event.EventMode,
		&event.EventType,
		&event.Status,
		&event.ContactEmail,
	)

	if err != nil {
		return entity.EventEntity{}, err
	}

	startDate, _ := pkg.ParseTime(*event.StartDate)
	endDate, _ := pkg.ParseTime(*event.EndDate)
	newEvent, insertionErr := e.eventRepo.CreateEvent(entity.EventEntity{
		Name:             event.Name,
		Description:      event.Description,
		SpectrumId:       event.SpectrumId,
		Status:           event.Status,
		ParticipantLimit: event.ParticipantLimit,
		StartDate:        &startDate,
		EndDate:          &endDate,
		EventMode:        event.EventMode,
		EventType:        event.EventType,
		IsFeatured:       event.IsFeatured,
		ContactEmail:     event.ContactEmail,
		VenueId:          event.VenueId,
	})
	if insertionErr != nil {
		return entity.EventEntity{}, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}
	return newEvent, nil

}

func (e EventUsecases) UpdateEvent(event entity.EventUpdateEntity) (entity.EventEntity, error) {

	emptyEvent := entity.EventEntity{}
	if !pkg.ValidateUUID(event.Id) {

		return emptyEvent, &customerrors.ValidationError{OrgError: "Invalid event id"}
	}

	validationErr := e.validateEventDetailsAndCheckExistence(
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
	)

	if validationErr != nil {
		return emptyEvent, validationErr
	}

	newEvent, updationErr := e.eventRepo.UpdateEvent(event.Id, event)

	if errors.Is(updationErr, gorm.ErrRecordNotFound) {
		return emptyEvent, &customerrors.NotFoundError{OrgError: "Event does not exist"}
	} else if updationErr != nil {
		return emptyEvent, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}

	return newEvent, nil
}

func (e EventUsecases) DeleteEvent(eventId string) error {

	if !pkg.ValidateUUID(eventId) {
		return &customerrors.ValidationError{OrgError: "Invalid event id"}
	}

	err := e.eventRepo.DeleteEvent(eventId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &customerrors.NotFoundError{OrgError: "Event does not exist"}
	}

	return nil
}

func (e EventUsecases) RetrieveEvents(limit int, page int) ([]entity.EventEntity, error) {

	offset := pkg.GetOffset(limit, page)
	events, err := e.eventRepo.RetrieveEvents(limit, offset)
	if err != nil {
		return events, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}

	return events, nil
}
func (e EventUsecases) validateEventDetailsAndCheckExistence(name *string,
	description *string,
	spectrumId *string,
	venueId *string,
	participantLimit *int,
	startDate *string,
	endDate *string,
	eventMode *enums.EventMode,
	eventType *enums.EventType,
	status *enums.EventStatus,
	contactEmail *string,

) error {

	if name != nil {
		nameErr := pkg.ValidateSpectrumOrEventName(*name)
		if nameErr != nil {
			return nameErr
		}
	}

	if description != nil {

		descriptionErr := pkg.ValidateSpectrumOrEventDescription(*description)
		if descriptionErr != nil {
			return descriptionErr
		}
	}

	if spectrumId != nil && !pkg.ValidateUUID(*spectrumId) {
		return &customerrors.ValidationError{OrgError: "Invalid spectrum id"}
	}

	if status != nil && !status.IsValid() {
		return &customerrors.ValidationError{OrgError: "Invalid event status"}

	}

	if participantLimit != nil && *(participantLimit) < 0 {
		return &customerrors.ValidationError{OrgError: "Invalid participant limit"}
	}

	var parsedStartTime *time.Time
	if startDate != nil {
		time, err := pkg.ParseTime(*startDate)
		if err != nil {
			return &customerrors.ValidationError{OrgError: "Invalid start date"}
		}
		parsedStartTime = &time

	}

	var parsedEndTime *time.Time
	if endDate != nil {
		time, err := pkg.ParseTime(*(endDate))
		if err != nil {
			return &customerrors.ValidationError{OrgError: "Invalid end date"}
		}
		parsedEndTime = &time
	}

	if parsedStartTime != nil {

		year, month, day := parsedStartTime.Date()
		fmt.Printf("Year : %d, month : %d, day : %d\n", year, month, day)
		now := time.Now()
		if year < now.Year() || month < now.Month() || day < now.Day() {
			return &customerrors.ValidationError{OrgError: "Start date should be today or after today"}
		} else if parsedEndTime != nil && parsedStartTime.After(*parsedEndTime) {
			return &customerrors.ValidationError{OrgError: "Start date should be on same day as end date or before end date"}
		}
	}

	if eventMode != nil && !eventMode.IsValid() {
		return &customerrors.ValidationError{OrgError: "Invalid event mode"}
	}

	if eventType != nil && !eventType.IsValid() {
		return &customerrors.ValidationError{OrgError: "Invalid event type"}
	}

	if contactEmail != nil && !pkg.ValidateEmail(*contactEmail) {
		return &customerrors.ValidationError{OrgError: "Invalid contact email"}
	}

	if venueId != nil && !pkg.ValidateUUID(*venueId) {
		return &customerrors.ValidationError{OrgError: "Invalid venue id"}
	}

	if spectrumId != nil {
		spectrumExists, spectrumExistenceErr := e.spectrumRepo.CheckSpectrumExists(*spectrumId)

		if spectrumExistenceErr != nil {
			return &customerrors.ServerError{OrgError: "Something went wrong while operating"}
		} else if !spectrumExists {
			return &customerrors.NotFoundError{OrgError: "Spectrum does not exist"}
		}
	}

	if venueId != nil {
		venueExists, venueExistanceErr := e.venueRepo.CheckVenueExists(*venueId)

		if venueExistanceErr != nil {
			return &customerrors.ServerError{OrgError: "Something went wrong while operating"}
		} else if !venueExists {
			return &customerrors.NotFoundError{OrgError: "Venue does not exist"}
		}

	}

	return nil
}
