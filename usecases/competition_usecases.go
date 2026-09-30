package usecases

import (
	customerrors "espectro/custom_errors"
	"espectro/entity"
	"espectro/enums"
	"espectro/pkg"
	"espectro/repository"
	"mime/multipart"
	"regexp"
	"time"

	"github.com/google/uuid"
)

type CompetitionUsecases struct {
	competitionRepo repository.CompetitionRepo
	venueRepo       repository.VenueRepo
	prizeRepo       repository.PrizeRepo
	mediaRepo       repository.MediaServiceRepo
}

func NewCompetitionUsecases(competitionRepo repository.CompetitionRepo, venueRepo repository.VenueRepo, prizeRepo repository.PrizeRepo, mediaRepo repository.MediaServiceRepo) CompetitionUsecases {
	return CompetitionUsecases{competitionRepo: competitionRepo, venueRepo: venueRepo, prizeRepo: prizeRepo, mediaRepo: mediaRepo}
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

	prizeExists, prizeCheckingErr := c.prizeRepo.PrizeExists(newCompetition.PrizeId)

	if prizeCheckingErr != nil {
		return empty, &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
	} else if !prizeExists {
		return empty, &customerrors.NotFoundError{DisplayError: "Prize does not exist"}
	}

	venueExists, venueCheckingErr := c.venueRepo.VenueExists(newCompetition.VenueId)

	if venueCheckingErr != nil {
		return empty, &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
	} else if !venueExists {
		return empty, &customerrors.NotFoundError{DisplayError: "Venue does not exist"}
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

	return nil
}
