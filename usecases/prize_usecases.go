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

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PrizeUsecases struct {
	prizeRepo repository.PrizeRepo
	mediaRepo repository.MediaServiceRepo
	eventRepo repository.EventRepo
}

func NewPrizeUsecases(prizeRepo repository.PrizeRepo, mediaRepo repository.MediaServiceRepo, eventRepo repository.EventRepo) PrizeUsecases {
	return PrizeUsecases{prizeRepo: prizeRepo, mediaRepo: mediaRepo, eventRepo: eventRepo}
}

func (p PrizeUsecases) CreatePrize(newPrize entity.PrizeCreateEntity) (entity.PrizeDBRetrieveEntity, error) {

	empty := entity.PrizeDBRetrieveEntity{}

	validationErr := p.validatePrizeData(nil, &newPrize.Title, &newPrize.Description, &newPrize.Amount, newPrize.Logo, &newPrize.EventId, &newPrize.Type)
	if validationErr != nil {
		return empty, validationErr
	}

	eventExists, eventCheckingErr := p.eventRepo.EventExists(newPrize.EventId)

	if eventCheckingErr != nil {
		return empty, &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
	} else if !eventExists {
		return empty, &customerrors.NotFoundError{DisplayError: "Event does not exist"}
	}

	var logoURL *string
	var logoPublicId *string

	prizeId := uuid.NewString()
	folderId := "prize/" + prizeId
	if newPrize.Logo != nil {

		url, pubId, err := p.mediaRepo.UploadFile(newPrize.Logo, folderId, false)
		if err != nil {
			return empty, &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
		}
		logoURL = &url
		logoPublicId = &pubId
	}

	createdPrize, creationErr := p.prizeRepo.CreatePrize(entity.PrizeDBCreateEntity{
		Id:          prizeId,
		Title:       newPrize.Title,
		Description: newPrize.Description,
		Amount:      newPrize.Amount,
		LogoURL:     logoURL,
		EventId:     newPrize.EventId,
		Type:        newPrize.Type,
	})

	if creationErr != nil {
		go func() {
			if logoPublicId != nil {
				p.mediaRepo.DeleteAssetsWithPublicIds([]string{*logoPublicId})
			}
		}()
		return empty, &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
	}

	return createdPrize, nil

}

func (p PrizeUsecases) UpdatePrize(prizeId string, newPrize entity.PrizeUpdateEntity) (entity.PrizeDBRetrieveEntity, error) {

	empty := entity.PrizeDBRetrieveEntity{}

	validationErr := p.validatePrizeData(&prizeId, newPrize.Title, newPrize.Description, newPrize.Amount, newPrize.Logo, newPrize.EventId, newPrize.Type)

	if validationErr != nil {
		return empty, validationErr
	}

	if newPrize.EventId != nil {

		exists, err := p.eventRepo.EventExists(*newPrize.EventId)
		if err != nil {
			return empty, &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
		}

		if !exists {
			return empty, &customerrors.NotFoundError{DisplayError: "Event does not exist"}
		}
	}

	var prevLogoPublicId *string
	var newLogoPublicId *string
	var newLogoURL *string

	folderId := "prize/" + prizeId
	if newPrize.Logo != nil {

		prevPubIds, prevPubIdsErr := p.mediaRepo.RetrieveAssetPublicIds(folderId, 1)

		if prevPubIdsErr != nil {
			return empty, &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
		}

		if len(prevPubIds) != 0 {
			prevLogoPublicId = &prevPubIds[0]
		}

		url, pubId, err := p.mediaRepo.UploadFile(newPrize.Logo, folderId, false)

		if err != nil {
			return empty, &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
		}

		newLogoURL = &url
		newLogoPublicId = &pubId
	}

	updatedPrize, updationErr := p.prizeRepo.UpdatePrize(prizeId, entity.PrizeDBUpdateEntity{
		Title:       newPrize.Title,
		Description: newPrize.Description,
		Amount:      newPrize.Amount,
		LogoURL:     newLogoURL,
		EventId:     newPrize.EventId,
		Type:        newPrize.Type,
	})

	if updationErr != nil {
		if newLogoPublicId != nil {
			go func() {
				p.mediaRepo.DeleteAssetsWithPublicIds([]string{*newLogoPublicId})
			}()
		}

		if errors.Is(updationErr, gorm.ErrRecordNotFound) {
			return empty, &customerrors.NotFoundError{DisplayError: "Prize does not exist"}
		} else {
			return empty, &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
		}

	} else {

		if prevLogoPublicId != nil {
			go func() {
				p.mediaRepo.DeleteAssetsWithPublicIds([]string{*prevLogoPublicId})
			}()
		}
		return updatedPrize, nil
	}

}

func (p PrizeUsecases) validatePrizeData(id *string, title *string, description *string, amount *float32, logo *multipart.FileHeader, eventId *string, prizeType *enums.PrizeTypeEnum) error {

	if id != nil && !pkg.ValidateUUID(*id) {
		return &customerrors.ValidationError{DisplayError: "Invalid id"}
	}

	if eventId != nil && !pkg.ValidateUUID(*eventId) {
		return &customerrors.ValidationError{DisplayError: "Invalid event id"}
	}

	if title != nil {
		if len(*title) > 100 {
			return &customerrors.ValidationError{DisplayError: "Title length should be less than or equal to 100"}
		} else if len(*title) < 5 {
			return &customerrors.ValidationError{DisplayError: "Title length should be greater than or equal to 5"}
		}

		regex := regexp.MustCompile(`[@!#$*^]`)
		contains := regex.MatchString(*title)
		if contains {
			return &customerrors.ValidationError{DisplayError: "Title should not contain any special character"}
		}
	}

	if description != nil {
		if len(*description) > 500 {
			return &customerrors.ValidationError{DisplayError: "Description length should be less than or equal to 500"}
		} else if len(*description) < 10 {
			return &customerrors.ValidationError{DisplayError: "Description length should be greater than or equal to 10"}
		}
	}

	if amount != nil && *amount < 0 {
		return &customerrors.ValidationError{DisplayError: "Amount should be positive"}
	}

	if prizeType != nil && !prizeType.IsValid() {
		return &customerrors.ValidationError{DisplayError: "Invalid type"}
	}

	if logo != nil {

		correct, imageErr := pkg.ValidateImage(logo)

		if imageErr != nil {
			return &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
		} else if !correct {
			return &customerrors.ValidationError{DisplayError: "Invalid image format"}
		}

		correctSize := pkg.ValidateImageSize(*logo)

		if !correctSize {
			return &customerrors.SizeError{DisplayError: "Logo size should be less than or equal to 2MB"}
		}
	}
	return nil
}
