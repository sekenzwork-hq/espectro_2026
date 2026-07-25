package usecases

import (
	"errors"
	customerrors "espectro/custom_errors"
	"espectro/database"
	"espectro/entity"
	"espectro/enums"
	"espectro/pkg"
	"espectro/repository"
	"time"

	"gorm.io/gorm"
)

type EventUsecases struct {
	eventRepo             repository.EventRepo
	spectrumRepo          repository.SpectrumRepo
	venueRepo             repository.VenueRepo
	transaction           database.TransactionManager
	eventsGalleryRepo     repository.EventGalleryRepo
	eventRegistrationRepo repository.EventRegistrationRepo
	userRepo              repository.UserRepository
}

func NewEventUsecases(
	eventRepo repository.EventRepo,
	spectrumRepo repository.SpectrumRepo,
	venueRepo repository.VenueRepo,
	transaction database.TransactionManager,
	eventsGalleryRepo repository.EventGalleryRepo,
	eventRegistrationRepo repository.EventRegistrationRepo,
	userRepo repository.UserRepository,
) EventUsecases {
	return EventUsecases{eventRepo: eventRepo,
		spectrumRepo:          spectrumRepo,
		venueRepo:             venueRepo,
		transaction:           transaction,
		eventsGalleryRepo:     eventsGalleryRepo,
		eventRegistrationRepo: eventRegistrationRepo,
		userRepo:              userRepo,
	}
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

	var startDate *time.Time
	var endDate *time.Time
	if event.StartDate != nil {
		date, _ := pkg.ParseTime(*event.StartDate)
		startDate = &date
	}
	if event.EndDate != nil {
		date, _ := pkg.ParseTime(*event.EndDate)
		endDate = &date
	}

	var newEvent entity.EventEntity

	transactionErr := e.transaction.Run(func() error {

		insertedEvent, insertionErr := e.eventRepo.CreateEvent(entity.EventEntity{
			Name:             event.Name,
			Description:      event.Description,
			SpectrumId:       event.SpectrumId,
			Status:           event.Status,
			ParticipantLimit: event.ParticipantLimit,
			StartDate:        startDate,
			EndDate:          endDate,
			EventMode:        event.EventMode,
			EventType:        event.EventType,
			IsFeatured:       event.IsFeatured,
			ContactEmail:     event.ContactEmail,
			VenueId:          event.VenueId,
		})
		if insertionErr != nil {
			return &customerrors.ServerError{OrgError: "Something went wrong while operating"}
		}

		incrementErr := e.spectrumRepo.IncrementTotalEventsCountBy1(event.SpectrumId)

		if incrementErr != nil {
			return &customerrors.ServerError{OrgError: "Something went wrong while operating"}
		}
		newEvent = insertedEvent
		return nil
	})

	return newEvent, transactionErr

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

	err := e.transaction.Run(func() error {

		spectrumId, err := e.eventRepo.RetrieveSpectrumId(eventId)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &customerrors.NotFoundError{OrgError: "Spectrum of event does not exist"}
		} else if err != nil {
			return &customerrors.ServerError{OrgError: "Something went wrong while operating"}
		}

		err = e.eventRepo.DeleteEvent(eventId)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &customerrors.NotFoundError{OrgError: "Event does not exist"}
		} else if err != nil {
			return &customerrors.ServerError{OrgError: "Something went wrong while operating"}
		}

		err = e.spectrumRepo.DecrementTotalEventsCountBy1(spectrumId)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &customerrors.NotFoundError{OrgError: "Spectrum of event does not exist"}
		} else if err != nil {
			return &customerrors.ServerError{OrgError: "Something went wrong while operating"}
		}

		err = e.eventsGalleryRepo.DeleteEventOrGallery(eventId)
		if err != nil {
			return &customerrors.ServerError{OrgError: "Something went wrong while operating"}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (e EventUsecases) RetrieveEvents(spectrumId *string, limit int, page int) ([]entity.EventEntity, error) {

	if spectrumId != nil && !pkg.ValidateUUID(*spectrumId) {
		return []entity.EventEntity{}, &customerrors.ValidationError{OrgError: "Invalid spectrum id"}
	}
	offset := pkg.GetOffset(limit, page)
	var events []entity.EventEntity
	var err error

	if spectrumId != nil {
		events, err = e.eventRepo.RetrieveEventsBySpectrumId(*spectrumId, limit, offset)
	} else {
		events, err = e.eventRepo.RetrieveEvents(limit, offset)
	}
	if err != nil {
		return events, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}

	return events, nil
}

func (e EventUsecases) Register(registrationDetails entity.EventRegistrationFromJsonEntity) (entity.EventRegistrationEntity, error) {

	empty := entity.EventRegistrationEntity{}
	if !pkg.ValidateUUID(registrationDetails.UserId) {
		return empty, &customerrors.ValidationError{OrgError: "Invalid user id"}
	} else if !pkg.ValidateUUID(registrationDetails.EventId) {
		return empty, &customerrors.ValidationError{OrgError: "Invalid event id"}
	}

	userExists, userErr := e.userRepo.UserExists(registrationDetails.UserId)
	if userErr != nil {
		return empty, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	} else if !userExists {
		return empty, &customerrors.AuthenticationError{OrgError: "User does not exist"}
	}

	eventExists, eventErr := e.eventRepo.EventExists(registrationDetails.EventId)
	if eventErr != nil {
		return empty, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	} else if !eventExists {
		return empty, &customerrors.NotFoundError{OrgError: "Event does not exist"}
	}

	alreadyRegistered, eventRegErr := e.eventRegistrationRepo.RegisterExists(registrationDetails.EventId, registrationDetails.UserId)
	if eventRegErr != nil {
		return empty, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	} else if !alreadyRegistered {
		return empty, &customerrors.ValidationError{OrgError: "User has been already registered"}
	}

	details, insertionErr := e.eventRegistrationRepo.Register(entity.EventRegistrationEntity{
		UserId:  registrationDetails.UserId,
		EventId: registrationDetails.EventId,
		Status:  enums.VerificationPending,
	})

	if insertionErr != nil {
		return empty, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}

	return details, nil

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
		time, err := pkg.ParseTime(*endDate)
		if err != nil {
			return &customerrors.ValidationError{OrgError: "Invalid end date"}
		}
		parsedEndTime = &time
	}

	if parsedStartTime != nil {
		year, month, day := parsedStartTime.Date()
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
