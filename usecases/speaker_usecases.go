package usecases

import (
	"errors"
	customerrors "espectro/custom_errors"
	"espectro/entity"
	"espectro/pkg"
	"espectro/repository"
	"mime/multipart"
	"regexp"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SpeakerUsecases struct {
	speakerRepo      repository.SpeakerRepo
	mediaRepo        repository.MediaServiceRepo
	eventSpeakerRepo repository.EventSpeakerRepo
	eventRepo        repository.EventRepo
}

func NewSpeakerUsecases(speakerRepo repository.SpeakerRepo, mediaRepo repository.MediaServiceRepo, eventRepo repository.EventRepo, eventSpeakerRepo repository.EventSpeakerRepo) SpeakerUsecases {
	return SpeakerUsecases{speakerRepo: speakerRepo, mediaRepo: mediaRepo, eventSpeakerRepo: eventSpeakerRepo, eventRepo: eventRepo}
}

func (s SpeakerUsecases) CreateSpeaker(newSpeaker entity.SpeakerCreateEntity) (entity.SpeakerDBRetrieveEntity, error) {

	empty := entity.SpeakerDBRetrieveEntity{}

	validationErr := s.validateSpeakerData(nil, &newSpeaker.Fullname, &newSpeaker.Bio, newSpeaker.ProfilePic, &newSpeaker.Country, &newSpeaker.PhoneNumber, &newSpeaker.Email)
	if validationErr != nil {
		return empty, validationErr
	}

	speakerId := uuid.NewString()
	folder := "speakers/" + speakerId
	var newProfilePicPublicId *string
	var newProfilePicURL *string

	if newSpeaker.ProfilePic != nil {

		url, pubId, err := s.mediaRepo.UploadFile(newSpeaker.ProfilePic, folder, false)
		if err != nil {
			return empty, &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
		}

		newProfilePicPublicId = &pubId
		newProfilePicURL = &url
	}

	createdSpeaker, creationErr := s.speakerRepo.CreateSpeaker(entity.SpeakerDBCreateEntity{
		Id:            speakerId,
		Fullname:      newSpeaker.Fullname,
		ProfilePicURL: newProfilePicURL,
		IsFeatured:    newSpeaker.IsFeatured,
		Bio:           newSpeaker.Bio,
		Country:       newSpeaker.Country,
		PhoneNumber:   newSpeaker.PhoneNumber,
		Email:         newSpeaker.Email,
	})

	if creationErr != nil {
		if newProfilePicPublicId != nil {
			go func() {
				s.mediaRepo.DeleteAssetsWithPublicIds([]string{*newProfilePicPublicId})
			}()
		}

		return empty, &customerrors.ServerError{DisplayError: "Seomething went wrong while operating"}
	}

	return createdSpeaker, nil
}

func (s SpeakerUsecases) UpdateSpeaker(speakerId string, newSpeaker entity.SpeakerUpdateEntity) (entity.SpeakerDBRetrieveEntity, error) {

	empty := entity.SpeakerDBRetrieveEntity{}

	validationErr := s.validateSpeakerData(&speakerId, newSpeaker.Fullname, newSpeaker.Bio, newSpeaker.ProfilePic, newSpeaker.Country, newSpeaker.PhoneNumber, newSpeaker.Email)

	if validationErr != nil {
		return empty, validationErr
	}

	folder := "speakers/" + speakerId
	var newProfilePicURL *string
	var newProfilePicPublicId *string
	var prevProfilePicPublicId *string

	if newSpeaker.ProfilePic != nil {

		prevPublicIds, prevPicRetrieveErr := s.mediaRepo.RetrieveAssetPublicIds(folder, 1)

		if prevPicRetrieveErr != nil {
			return empty, &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
		}

		if len(prevPublicIds) != 0 {
			prevProfilePicPublicId = &prevPublicIds[0]
		}

		url, pubId, err := s.mediaRepo.UploadFile(newSpeaker.ProfilePic, folder, false)

		if err != nil {
			return empty, &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
		}

		newProfilePicPublicId = &pubId
		newProfilePicURL = &url
	}

	updatedSpeaker, updationErr := s.speakerRepo.UpdateSpeaker(speakerId, entity.SpeakerDBUpdateEntity{
		Fullname:      newSpeaker.Fullname,
		ProfilePicURL: newProfilePicURL,
		IsFeatured:    newSpeaker.IsFeatured,
		Bio:           newSpeaker.Bio,
		Country:       newSpeaker.Country,
		PhoneNumber:   newSpeaker.PhoneNumber,
		Email:         newSpeaker.Email,
	})

	if updationErr != nil {

		if newProfilePicPublicId != nil {
			go func() {
				s.mediaRepo.DeleteAssetsWithPublicIds([]string{*newProfilePicPublicId})
			}()
		}

		if errors.Is(updationErr, gorm.ErrRecordNotFound) {
			return empty, &customerrors.NotFoundError{DisplayError: "Speaker does not exist"}
		}
		return empty, &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
	}

	if prevProfilePicPublicId != nil {
		go func() {
			s.mediaRepo.DeleteAssetsWithPublicIds([]string{*prevProfilePicPublicId})
		}()
	}

	return updatedSpeaker, nil

}

func (s SpeakerUsecases) DeleteSpeaker(speakerId string) error {

	if !pkg.ValidateUUID(speakerId) {
		return &customerrors.ValidationError{DisplayError: "Invalid speaker id"}
	}

	err := s.speakerRepo.DeleteSpeaker(speakerId)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &customerrors.NotFoundError{DisplayError: "Speaker does not exist"}
	} else if err != nil {
		return &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
	}

	go func() {
		s.mediaRepo.DeleteFile("speakers/", speakerId)
	}()

	return nil
}

func (s SpeakerUsecases) AddSpeakerToEvent(eventId string, speakerId string) error {

	if !pkg.ValidateUUID(eventId) {
		return &customerrors.ValidationError{DisplayError: "Invalid event id"}
	}

	if !pkg.ValidateUUID(speakerId) {
		return &customerrors.ValidationError{DisplayError: "Invalid speaker id"}
	}

	speakerExists, speakerErr := s.speakerRepo.SpeakerExists(speakerId)

	if speakerErr != nil {
		return &customerrors.ServerError{DisplayError: "Something went wrong while operatinsg"}
	} else if !speakerExists {
		return &customerrors.NotFoundError{DisplayError: "Speaker does not exist"}
	}
	eventExists, eventErr := s.eventRepo.EventExists(eventId)

	if eventErr != nil {
		return &customerrors.ServerError{DisplayError: "Something went wrong while operatinsg"}
	} else if !eventExists {
		return &customerrors.NotFoundError{DisplayError: "Event does not exist"}
	}

	isAlreadyAdded, checkingErr := s.eventSpeakerRepo.SpeakerIsAlreadyAddedToEvent(speakerId, eventId)

	if checkingErr != nil {
		return &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
	} else if isAlreadyAdded {
		return &customerrors.ValidationError{DisplayError: "Speaker is already added to the event"}
	}

	err := s.eventSpeakerRepo.AddSpeakerToEvent(speakerId, eventId)

	if err != nil {
		return &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
	}

	return nil
}

func (s SpeakerUsecases) RemoveSpeakerFromEvent(speakerId string, eventId string) error {

	if !pkg.ValidateUUID(speakerId) {
		return &customerrors.ValidationError{DisplayError: "Invalid speaker id"}
	}

	if !pkg.ValidateUUID(eventId) {
		return &customerrors.ValidationError{DisplayError: "Invalid event id"}
	}

	err := s.eventSpeakerRepo.RemoveSpeakerFromEvent(speakerId, eventId)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &customerrors.NotFoundError{DisplayError: "Speaker is not added to the event"}
	} else if err != nil {
		return &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
	}

	return nil
}

func (s SpeakerUsecases) RetrieveSpeaker(eventId string, page int) ([]entity.SpeakerDBRetrieveEntity, error) {

	empty := []entity.SpeakerDBRetrieveEntity{}

	if !pkg.ValidateUUID(eventId) {
		return empty, &customerrors.ValidationError{DisplayError: "Invalid event id"}
	}

	offset := pkg.GetOffset(50, page)
	speakers, err := s.speakerRepo.RetrieveSpeaker(eventId, offset)

	if err != nil {
		return empty, &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
	}

	return speakers, nil
}

func (s SpeakerUsecases) validateSpeakerData(id *string, fullname *string, bio *string, profilePic *multipart.FileHeader, country *string, phoneNumber *string, email *string) error {

	if id != nil && !pkg.ValidateUUID(*id) {
		return &customerrors.ValidationError{DisplayError: "Invalid speaker id"}
	}

	if fullname != nil {

		err := pkg.ValidateFullname(*fullname)
		if err != nil {
			return err
		}
	}

	if bio != nil {

		if len(*bio) < 5 {
			return &customerrors.ValidationError{DisplayError: "Bio length should be greater than or equal to 5"}
		} else if len(*bio) > 100 {
			return &customerrors.ValidationError{DisplayError: "Bio length should be less than or equal to 100"}
		}

		regex := regexp.MustCompile(`[@#$*&^!;+=\-]`)
		contains := regex.MatchString(*bio)
		if contains {
			return &customerrors.ValidationError{DisplayError: "Bio should not contain any special character"}
		}
	}

	if profilePic != nil {

		correct, imageCorrectErr := pkg.ValidateImage(profilePic)

		if imageCorrectErr != nil {
			return &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
		} else if !correct {
			return &customerrors.ValidationError{DisplayError: "Invalid image format"}
		}

		sizeCorrect := pkg.ValidateImageSize(*profilePic)

		if !sizeCorrect {
			return &customerrors.ServerError{DisplayError: "Profile pic size should be less than or equal to 2MB"}
		}
	}

	if country != nil && !pkg.ValidateCountryOrState(*country) {
		return &customerrors.ValidationError{DisplayError: "Invalid country"}
	}

	if phoneNumber != nil && !pkg.ValidatePhoneNumber(*phoneNumber) {
		return &customerrors.ValidationError{DisplayError: "Invalid phone number"}
	}

	if email != nil && !pkg.ValidateEmail(*email) {
		return &customerrors.ValidationError{DisplayError: "Invalid email"}
	}

	return nil

}
