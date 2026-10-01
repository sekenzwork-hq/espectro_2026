package usecases

import (
	"errors"
	customerrors "espectro/custom_errors"
	"espectro/entity"
	"espectro/enums"
	"espectro/pkg"
	"espectro/repository"
	"mime/multipart"
	"regexp"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CompetitionUsecases struct {
	competitionRepo repository.CompetitionRepo
	venueRepo       repository.VenueRepo
	prizeRepo       repository.PrizeRepo
	eventRepo       repository.EventRepo
	mediaRepo       repository.MediaServiceRepo
}

func NewCompetitionUsecases(competitionRepo repository.CompetitionRepo, venueRepo repository.VenueRepo, prizeRepo repository.PrizeRepo, eventRepo repository.EventRepo, mediaRepo repository.MediaServiceRepo) CompetitionUsecases {
	return CompetitionUsecases{competitionRepo: competitionRepo, venueRepo: venueRepo, prizeRepo: prizeRepo, mediaRepo: mediaRepo, eventRepo: eventRepo}
}

func (c CompetitionUsecases) CreateCompetition(newCompetition entity.CompetitionCreateEntity) (entity.CompetitionDBRetrieveEntity, error) {

	empty := entity.CompetitionDBRetrieveEntity{}

	validationErr := c.validateCompetitionData(
		nil,
		&newCompetition.EventId,
		&newCompetition.Title,
		&newCompetition.Description,
		&newCompetition.Rules,
		&newCompetition.Status,
		&newCompetition.OpeningTime,
		&newCompetition.ClosingTime,
		&newCompetition.ResultTime,
		&newCompetition.PrizeId,
		&newCompetition.VenueId,
		&newCompetition.RegistrationFee,
		newCompetition.Logo,
	)

	if validationErr != nil {
		return empty, validationErr
	}

	var newLogoURL *string
	var newLogoPublicId *string

	competitionId := uuid.NewString()
	folder := "competitions/" + competitionId

	if newCompetition.Logo != nil {

		url, pubId, err := c.mediaRepo.UploadFile(newCompetition.Logo, folder, false)

		if err != nil {
			return empty, &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
		}
		newLogoURL = &url
		newLogoPublicId = &pubId
	}

	createdCompetition, creationErr := c.competitionRepo.CreateCompetition(entity.CompetitionDBCreateEntity{
		Id:              competitionId,
		EventId:         newCompetition.EventId,
		Title:           newCompetition.Title,
		Description:     newCompetition.Description,
		Rules:           newCompetition.Rules,
		Status:          newCompetition.Status,
		OpeningTime:     newCompetition.OpeningTime,
		ClosingTime:     newCompetition.ClosingTime,
		ResultTime:      newCompetition.ResultTime,
		RegistrationFee: newCompetition.RegistrationFee,
		PrizeId:         newCompetition.PrizeId,
		VenueId:         newCompetition.VenueId,
		LogoURL:         newLogoURL,
	})

	if creationErr != nil {
		if newLogoPublicId != nil {
			go func() {
				c.mediaRepo.DeleteAssetsWithPublicIds([]string{*newLogoPublicId})
			}()
		}

		return empty, &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
	}

	return createdCompetition, nil

}

func (c CompetitionUsecases) UpdateCompetition(competitionId string, newCompetition entity.CompetitionUpdateEntity) (entity.CompetitionDBRetrieveEntity, error) {

	var empty entity.CompetitionDBRetrieveEntity

	validationErr := c.validateCompetitionData(
		&competitionId,
		newCompetition.EventId,
		newCompetition.Title,
		newCompetition.Description,
		newCompetition.Rules,
		newCompetition.Status,
		newCompetition.OpeningTime,
		newCompetition.ClosingTime,
		newCompetition.ResultTime,
		newCompetition.PrizeId,
		newCompetition.VenueId,
		newCompetition.RegistrationFee,
		newCompetition.Logo,
	)

	if validationErr != nil {
		return empty, validationErr
	}

	var newLogoURL *string
	var newLogoPublicId *string
	var prevLogoPubId *string

	folder := "competitions/" + competitionId

	if newCompetition.Logo != nil {

		prevPubIds, prevErr := c.mediaRepo.RetrieveAssetPublicIds(folder, 1)

		if prevErr != nil {
			return empty, &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
		}

		if len(prevPubIds) != 0 {
			prevLogoPubId = &prevPubIds[0]
		}

		url, pubId, err := c.mediaRepo.UploadFile(newCompetition.Logo, folder, false)

		if err != nil {
			return empty, &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
		}
		newLogoURL = &url
		newLogoPublicId = &pubId
	}

	updatedCompetition, updationErr := c.competitionRepo.UpdateCompetition(competitionId, entity.CompetitionDBUpdateEntity{
		EventId:         newCompetition.EventId,
		Title:           newCompetition.Title,
		Description:     newCompetition.Description,
		Rules:           newCompetition.Rules,
		Status:          newCompetition.Status,
		OpeningTime:     newCompetition.OpeningTime,
		ClosingTime:     newCompetition.ClosingTime,
		ResultTime:      newCompetition.ResultTime,
		RegistrationFee: newCompetition.RegistrationFee,
		PrizeId:         newCompetition.PrizeId,
		VenueId:         newCompetition.VenueId,
		LogoURL:         newLogoURL,
	})

	if updationErr != nil {
		if newLogoPublicId != nil {
			go func() {
				c.mediaRepo.DeleteAssetsWithPublicIds([]string{*newLogoPublicId})
			}()
		}

		if errors.Is(updationErr, gorm.ErrRecordNotFound) {
			return empty, &customerrors.NotFoundError{DisplayError: "Competition does not exist"}
		} else {
			return empty, &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
		}
	} else if prevLogoPubId != nil {
		go func() {
			c.mediaRepo.DeleteAssetsWithPublicIds([]string{*prevLogoPubId})
		}()
	}

	return updatedCompetition, nil
}

func (c CompetitionUsecases) DeleteCompetition(competitionId string) error {

	if !pkg.ValidateUUID(competitionId) {
		return &customerrors.ValidationError{DisplayError: "Invalid competition id"}
	}

	err := c.competitionRepo.DeleteCompetition(competitionId)

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &customerrors.NotFoundError{DisplayError: "Competition does not exist"}
		} else {
			return &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
		}
	} else {
		go func() {
			c.mediaRepo.DeleteFile("competitions/", competitionId)
		}()
	}

	return nil
}

func (c CompetitionUsecases) RetrieveCompetitions(eventId string, page int) ([]entity.CompetitionDBRetrieveEntity, error) {

	if !pkg.ValidateUUID(eventId) {
		return []entity.CompetitionDBRetrieveEntity{}, &customerrors.ValidationError{DisplayError: "Invalid event id"}
	}

	offset := pkg.GetOffset(50, page)

	competitions, err := c.competitionRepo.RetrieveCompetitions(eventId, offset)

	if err != nil {
		return competitions, &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
	}

	return competitions, nil
}

func (c CompetitionUsecases) validateCompetitionData(id *string, eventId *string, title *string, description *string, rules *string, status *enums.CompetitionStatusEnum, openingTime *time.Time, closingTime *time.Time, resultTime *time.Time, prizeId *string, venueId *string, registrationFee *float32, logo *multipart.FileHeader) error {

	if id != nil && !pkg.ValidateUUID(*id) {
		return &customerrors.ValidationError{DisplayError: "Invalid competition id"}
	}

	if eventId != nil && !pkg.ValidateUUID(*eventId) {
		return &customerrors.ValidationError{DisplayError: "Invalid event id"}
	}

	if title != nil {

		if len(*title) < 5 {
			return &customerrors.ValidationError{DisplayError: "Title length should be greater than or equal to 5"}
		} else if len(*title) > 100 {
			return &customerrors.ValidationError{DisplayError: "Title length should be less than or equal to 100"}
		}

		regex := regexp.MustCompile(`[!@#$*\-+=;\{\}]`)

		if regex.MatchString(*title) {
			return &customerrors.ValidationError{DisplayError: "Title should not contain any special character"}
		}
	}

	if description != nil {

		if len(*description) < 5 {
			return &customerrors.ValidationError{DisplayError: "Description length should be greater than or equal to 5"}
		} else if len(*description) > 300 {
			return &customerrors.ValidationError{DisplayError: "Description length should be less than or equal to 300"}
		}
	}

	if rules != nil {

		if len(*rules) < 5 {
			return &customerrors.ValidationError{DisplayError: "Rules length should be less than or equal to 5"}
		} else if len(*rules) > 1000 {
			return &customerrors.ValidationError{DisplayError: "Rules length should be less than or equal to 1000"}
		}
	}

	if status != nil && !status.IsValid() {
		return &customerrors.ValidationError{DisplayError: "Invalid status"}
	}

	if registrationFee != nil {
		value := *registrationFee
		if value < 0 {
			return &customerrors.ValidationError{DisplayError: "Registration fee should be positive"}
		}
	}

	now := time.Now()

	if openingTime != nil && now.After(*openingTime) {
		return &customerrors.ValidationError{DisplayError: "Opening time must be past the current time"}
	}

	if closingTime != nil && now.After(*closingTime) {
		return &customerrors.ValidationError{DisplayError: "Closing time must be past the current time"}
	}

	if resultTime != nil && now.After(*resultTime) {
		return &customerrors.ValidationError{DisplayError: "Result time must be past the current time"}
	}

	if prizeId != nil && !pkg.ValidateUUID(*prizeId) {
		return &customerrors.ValidationError{DisplayError: "Invalid prize id"}
	}
	if venueId != nil && !pkg.ValidateUUID(*venueId) {
		return &customerrors.ValidationError{DisplayError: "Invalid venue id"}
	}

	if logo != nil {

		isCorrect, err := pkg.ValidateImage(logo)

		if err != nil {
			return &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
		} else if !isCorrect {
			return &customerrors.ValidationError{DisplayError: "Invalid image format"}
		}

		isCorrectSize := pkg.ValidateImageSize(*logo)

		if !isCorrectSize {
			return &customerrors.ValidationError{DisplayError: "Image size should be less than or equal to 2MB"}
		}
	}

	if eventId != nil {
		eventExists, eventCheckingErr := c.eventRepo.EventExists(*eventId)

		if eventCheckingErr != nil {
			return &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
		} else if !eventExists {
			return &customerrors.NotFoundError{DisplayError: "Event does not exist"}
		}
	}
	if prizeId != nil {

		prizeExists, prizeCheckingErr := c.prizeRepo.PrizeExists(*prizeId)

		if prizeCheckingErr != nil {
			return &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
		} else if !prizeExists {
			return &customerrors.NotFoundError{DisplayError: "Prize does not exist"}
		}
	}

	if venueId != nil {
		venueExists, venueCheckingErr := c.venueRepo.VenueExists(*venueId)

		if venueCheckingErr != nil {
			return &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
		} else if !venueExists {
			return &customerrors.NotFoundError{DisplayError: "Venue does not exist"}
		}
	}

	return nil
}
