package usecases

import (
	"errors"
	customerrors "espectro/custom_errors"
	"espectro/database"
	"espectro/entity"
	"espectro/enums"
	"espectro/models"
	"espectro/pkg"
	"espectro/repository"
	"mime/multipart"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SpectrumUsecases struct {
	spectrumRepo repository.SpectrumRepo
	mediaRepo    repository.MediaServiceRepo
	eventRepo    repository.EventRepo
	transaction  database.TransactionManager
}

func NewSpectrumUsecases(
	spectrumRepo repository.SpectrumRepo,
	mediaRepo repository.MediaServiceRepo,
	eventRepo repository.EventRepo,
	transaction database.TransactionManager,
) SpectrumUsecases {
	return SpectrumUsecases{spectrumRepo: spectrumRepo, mediaRepo: mediaRepo, eventRepo: eventRepo, transaction: transaction}
}

func (s SpectrumUsecases) CreateSpectrum(
	name string,
	shortDescription string,
	description string,
	status enums.SpectrumStatus,
	imageFiles []*multipart.FileHeader,
	videoFile *multipart.FileHeader,
	logoFile *multipart.FileHeader,

) (entity.SpectrumEntity, error) {

	emptySpectrum := entity.SpectrumEntity{}

	validationErr := s.validateSpectrumData(&name, &shortDescription, &description, &status, imageFiles, videoFile, logoFile)
	if validationErr != nil {
		return emptySpectrum, validationErr
	}

	spectrumId := uuid.NewString()

	media, mediaErr := s.uploadMediaForSpectrum(spectrumId, logoFile, videoFile, imageFiles)
	if mediaErr != nil {
		return emptySpectrum, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}

	spectrumData := entity.SpectrumEntity{
		Id:               spectrumId,
		Name:             name,
		ShortDescription: shortDescription,
		Description:      description,
		Status:           enums.PendingSpectrum,
		ImageUrls:        media.ImageUrls,
		LogoUrl:          media.LogoUrl,
		VideoUrl:         media.VideoUrl,
		TotalEvents:      0,
	}

	spectrum, creationErr := s.spectrumRepo.CreateSpectrum(spectrumData)

	if creationErr != nil {
		go s.mediaRepo.DeleteMutipleFiles("spectrum/"+spectrumId, []string{"logo", "video", "images"})
	}

	if creationErr != nil {
		return emptySpectrum, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}

	return spectrum, nil
}

func (s SpectrumUsecases) UpdateSpectrum(
	spectrumId string,
	name *string,
	shortDescription *string,
	description *string,
	status *enums.SpectrumStatus,
	imageFiles []*multipart.FileHeader,
	videoFile *multipart.FileHeader,
	logoFile *multipart.FileHeader,
) (entity.SpectrumEntity, error) {

	emptySpectrum := entity.SpectrumEntity{}
	serverErr := &customerrors.ServerError{OrgError: "Something went wrong while operating"}

	if correct := pkg.ValidateUUID(spectrumId); !correct {
		return emptySpectrum, &customerrors.ValidationError{OrgError: "Invalid spectrum id"}
	}

	validationErr := s.validateSpectrumData(name, shortDescription, description, status, imageFiles, videoFile, logoFile)
	if validationErr != nil {
		return emptySpectrum, validationErr
	}

	media, mediaErr := s.uploadMediaForSpectrum(spectrumId, logoFile, videoFile, imageFiles)
	if mediaErr != nil {
		return emptySpectrum, serverErr
	}

	newSpectrum, updationErr := s.spectrumRepo.UpdateSpectrum(spectrumId,
		name,
		shortDescription,
		description,
		status,
		media.LogoUrl,
		media.VideoUrl,
		media.ImageUrls,
	)

	if updationErr != nil {
		go s.mediaRepo.DeleteAssetsWithPublicIds(media.NewImagePublicIds)
	}
	if errors.Is(updationErr, gorm.ErrRecordNotFound) {
		return emptySpectrum, &customerrors.NotFoundError{OrgError: "Spectrum does not exist"}
	} else if updationErr != nil {
		return emptySpectrum, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}

	go s.mediaRepo.DeleteAssetsWithPublicIds(media.PreviousImagePublicIds)

	return newSpectrum, nil

}

func (s SpectrumUsecases) DeleteSpectrum(spectrumId string) error {

	isIdCorrect := pkg.ValidateUUID(spectrumId)

	if !isIdCorrect {
		return &customerrors.ValidationError{OrgError: "Invalid spectrum id"}
	}

	err := s.transaction.Run(func() error {
		err := s.spectrumRepo.DeleteSpectrum(spectrumId)

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &customerrors.NotFoundError{OrgError: "Spectrum does not exist"}
		} else if err != nil {
			return &customerrors.ServerError{OrgError: "Something went wrong while operating"}
		}

		err = s.eventRepo.DeleteEventBySpectrumId(spectrumId)
		if err != nil {
			return &customerrors.ServerError{OrgError: "Something went wrong while operating"}
		}
		return nil
	})

	return err
}

func (s SpectrumUsecases) RetrieveSpectrums(limit int, page int) ([]entity.SpectrumEntity, error) {

	emptySpectrums := []entity.SpectrumEntity{}

	if limit > 150 {
		return emptySpectrums, &customerrors.SizeError{OrgError: "Limit should be less than or equal to 150"}
	}

	offset := pkg.GetOffset(limit, page)
	spectrums, spectrumsErr := s.spectrumRepo.RetrieveSpectrums(offset, limit)
	if spectrumsErr != nil {
		return emptySpectrums, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}
	return spectrums, nil
}

func (s SpectrumUsecases) validateSpectrumData(
	name *string,
	shortDescription *string,
	description *string,
	status *enums.SpectrumStatus,
	imageFiles []*multipart.FileHeader,
	videoFile *multipart.FileHeader,
	logoFile *multipart.FileHeader,
) error {

	if name != nil {
		nameErr := pkg.ValidateSpectrumOrEventName(*(name))
		if nameErr != nil {
			return nameErr
		}
	}

	if shortDescription != nil {
		shortDescriptionErr := pkg.ValidateSpectrumShortDescription(*(shortDescription))
		if shortDescriptionErr != nil {
			return shortDescriptionErr
		}
	}

	if description != nil {
		descriptionErr := pkg.ValidateSpectrumOrEventDescription(*(description))
		if descriptionErr != nil {
			return descriptionErr
		}

	}

	if status != nil {
		if !enums.SpectrumStatus(*status).IsValid() {
			return &customerrors.ValidationError{OrgError: "Invalid status"}
		}
	}

	if len(imageFiles) > 10 {
		return &customerrors.ValidationError{OrgError: "Maximum number of images is 10"}
	}

	if logoFile != nil && pkg.ValidateImageSize(*logoFile) {
		return &customerrors.SizeError{OrgError: "Logo image size should be less than or equal to 2 MB"}
	} else if videoFile != nil && pkg.ValidateVideoSize(*videoFile) {
		return &customerrors.SizeError{OrgError: "Video size should be less than or equal to 50 MB"}
	}

	for i := range imageFiles {
		if imageFiles[i] != nil && !pkg.ValidateImageSize(*imageFiles[i]) {
			return &customerrors.SizeError{OrgError: "Size of each image should be less than or equal to 2 MB"}
		}

		valid, err := pkg.ValidateImage(imageFiles[i])
		if err != nil {
			return &customerrors.ServerError{OrgError: "Something went wrong"}
		} else if !valid {
			return &customerrors.ValidationError{OrgError: "Invalid image format"}
		}
	}

	return nil

}

func (s SpectrumUsecases) uploadMediaForSpectrum(
	spectrumId string,
	logoFile *multipart.FileHeader,
	videoFile *multipart.FileHeader,
	imageFiles []*multipart.FileHeader,
) (models.SpectrumMediaModel, error) {

	var logoUrl *string
	var videoUrl *string
	var imageUrls []string

	baseFolder := "spectrum/" + spectrumId
	logoFolder := baseFolder + "/logo"
	imagesFolder := baseFolder + "/images"
	videoFolder := baseFolder + "/video"
	emptyModel := models.SpectrumMediaModel{}

	funcToDeleteUploadedMedia := func() {
		s.mediaRepo.DeleteMutipleFiles(baseFolder, []string{logoFolder, videoFolder})
	}

	if logoFile != nil {
		url, _, urlErr := s.mediaRepo.UploadFile(logoFile, logoFolder, true)
		if urlErr != nil {
			return emptyModel, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
		}
		logoUrl = &url
	}

	if videoFile != nil {
		url, _, urlErr := s.mediaRepo.UploadFile(logoFile, videoFolder, true)
		if urlErr != nil {
			//Deleting the above uploaded logo (if it is provided ) if video is failed while uploading
			if logoUrl != nil {
				go s.mediaRepo.DeleteFile(baseFolder, logoFolder)
			}
			return emptyModel, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
		}
		videoUrl = &url
	}

	var prevImagePublicIds []string
	var newImagePublicIds []string
	if len(imageFiles) != 0 {
		prevPublicIds, err := s.mediaRepo.RetrieveAssetPublicIds(imagesFolder)
		if err != nil {
			funcToDeleteUploadedMedia()
			return emptyModel, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
		}

		urls, newPublicIds, urlsErr := s.mediaRepo.UploadFiles(imageFiles, imagesFolder, false)
		if urlsErr != nil {
			//Deleting the above upload video and logo (if those are provided) and deleting rest of the images
			//we were uploading if it is failed.
			funcToDeleteUploadedMedia()
			return emptyModel, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
		}
		imageUrls = urls
		prevImagePublicIds = prevPublicIds
		newImagePublicIds = newPublicIds
	}

	return models.SpectrumMediaModel{LogoUrl: logoUrl, VideoUrl: videoUrl, ImageUrls: imageUrls, PreviousImagePublicIds: prevImagePublicIds, NewImagePublicIds: newImagePublicIds}, nil

}
