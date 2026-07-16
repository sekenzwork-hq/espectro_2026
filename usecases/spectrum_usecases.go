package usecases

import (
	"errors"
	customerrors "espectro/custom_errors"
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
	repo      repository.SpectrumRepo
	mediaRepo repository.MediaServiceRepo
}

func NewSpectrumUsecases(repo repository.SpectrumRepo, mediaRepo repository.MediaServiceRepo) SpectrumUsecases {
	return SpectrumUsecases{repo: repo, mediaRepo: mediaRepo}
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

	spectrumId, idErr := uuid.NewUUID()

	if idErr != nil {
		return emptySpectrum, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}

	media, mediaErr := s.uploadMediaForSpectrum(spectrumId.String(), logoFile, videoFile, imageFiles)

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

	spectrum, creationErr := s.repo.CreateSpectrum(spectrumData)

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

	newSpectrum, updationErr := s.repo.UpdateSpectrum(spectrumId,
		name,
		shortDescription,
		description,
		status,
		media.LogoUrl,
		media.VideoUrl,
		media.ImageUrls,
	)

	if errors.Is(updationErr, gorm.ErrRecordNotFound) {
		return emptySpectrum, &customerrors.NotFoundError{OrgError: "Spectrum does not exist"}
	} else if updationErr != nil {
		return emptySpectrum, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}

	return newSpectrum, nil

}

func (s SpectrumUsecases) DeleteSpectrum(spectrumId string) error {

	isIdCorrect := pkg.ValidateUUID(spectrumId)

	if !isIdCorrect {
		return &customerrors.ValidationError{OrgError: "Invalid spectrum id"}
	}

	deletionErr := s.repo.DeleteSpectrum(spectrumId)

	if errors.Is(deletionErr, gorm.ErrRecordNotFound) {
		return &customerrors.NotFoundError{OrgError: "Spectrum does not exist"}
	} else if deletionErr != nil {
		return &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}

	return nil
}

func (s SpectrumUsecases) RetrieveSpectrums(limit int, page int) ([]entity.SpectrumEntity, error) {

	emptySpectrums := []entity.SpectrumEntity{}

	if limit > 150 {
		return emptySpectrums, &customerrors.SizeError{OrgError: "Limit should be less than or equal to 150"}
	}

	offset := pkg.GetOffset(limit, page)
	spectrums, spectrumsErr := s.repo.RetrieveSpectrums(offset, limit)
	if spectrumsErr != nil {
		return emptySpectrums, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}
	return spectrums, nil
}

/*
For validating the basic spectrum data and if one of those is null, then it won't validate that.
Maybe this function is called for updating few fields only. Here passing the media files only for checking size and limit.
*/
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
		return &customerrors.SizeError{OrgError: "Maximum number of images is 10"}
	}

	if logoFile != nil && pkg.BytesToMB(logoFile.Size) > 2 {
		return &customerrors.SizeError{OrgError: "Logo image size should be less than or equal to 2 MB"}
	} else if videoFile != nil && pkg.BytesToMB(videoFile.Size) > 50 {
		return &customerrors.SizeError{OrgError: "Video size should be less than or equal to 50 MB"}
	}

	for i := range imageFiles {
		mb := pkg.BytesToMB(imageFiles[i].Size)
		if mb > 2 {
			return &customerrors.SizeError{OrgError: "Images size should be less than or equal to 2 MB"}
		}
	}

	return nil

}

/*
For uploading the spectrum media such as video, logo and other images len(10) and also deleting-
all those media we upload if one of them fail to upload
*/
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

	if logoFile != nil {
		url, urlErr := s.mediaRepo.UploadFile(logoFile, logoFolder)
		if urlErr != nil {
			return emptyModel, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
		}
		logoUrl = &url
	}

	if videoFile != nil {
		url, urlErr := s.mediaRepo.UploadFile(logoFile, videoFolder)
		if urlErr != nil {
			//Deleting the above uploaded logo (if it is provided ) if video is failed while uploading
			if logoUrl != nil {
				s.mediaRepo.DeleteFolderWithFiles(baseFolder, logoFolder)
			}
			return emptyModel, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
		}
		videoUrl = &url
	}

	if len(imageFiles) != 0 {
		urls, urlsErr := s.mediaRepo.UploadFiles(imageFiles, imagesFolder)
		if urlsErr != nil {
			//Deleting the above upload video and logo (if those are provided) and deleting rest of the images
			//we were uploading if it is failed.
			s.mediaRepo.DeleteMutiFoldersWithFiles(baseFolder, []string{logoFolder, videoFolder, imagesFolder})

			return emptyModel, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
		}
		imageUrls = urls

	}

	return models.SpectrumMediaModel{LogoUrl: logoUrl, VideoUrl: videoUrl, ImageUrls: imageUrls}, nil

}
