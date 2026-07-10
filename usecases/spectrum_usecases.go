package usecases

import (
	customerrors "espectro/custom_errors"
	"espectro/entity"
	"espectro/enums"
	"espectro/pkg"
	"espectro/repository"
	"mime/multipart"

	"github.com/google/uuid"
)

type SpectrumUsecases struct {
	repo repository.SpectrumRepo
}

func NewSpectrumUsecases(repo repository.SpectrumRepo) SpectrumUsecases {
	return SpectrumUsecases{repo: repo}
}

func (s SpectrumUsecases) CreateSpectrum(
	name string,
	shortDescription string,
	description string,
	status enums.SpectrumStatus,
	imageFiles []*multipart.FileHeader,
	videoFile *multipart.FileHeader,
	logoFile *multipart.FileHeader,

) (entity.Spectrum, error) {

	nameErr := pkg.ValidateSpectrumOrEventName(name)

	if nameErr != nil {
		return entity.Spectrum{}, nameErr
	}

	shortDescriptionErr := pkg.ValidateSpectrumShortDescription(shortDescription)

	if shortDescriptionErr != nil {
		return entity.Spectrum{}, shortDescriptionErr
	}

	descriptionErr := pkg.ValidateSpectrumOrEventDescription(description)

	if descriptionErr != nil {
		return entity.Spectrum{}, descriptionErr
	}

	if len(imageFiles) > 10 {
		return entity.Spectrum{}, &customerrors.SizeError{OrgError: "Maximum number of images is 10"}
	}

	var logoUrl *string
	var videoUrl *string
	var imageUrls []string

	if pkg.BytesToMB(logoFile.Size) > 2 {
		return entity.Spectrum{}, &customerrors.SizeError{OrgError: "Logo image size should be less than or equal to 2 MB"}
	} else if pkg.BytesToMB(videoFile.Size) > 50 {
		return entity.Spectrum{}, &customerrors.SizeError{OrgError: "Video size should be less than or equal to 50 MB"}
	}

	for i := range imageFiles {
		mb := pkg.BytesToMB(imageFiles[i].Size)

		if mb > 2 {
			return entity.Spectrum{}, &customerrors.SizeError{OrgError: "Images size should be less than or equal to 2 MB"}
		}
	}

	spectrumId, idErr := uuid.NewUUID()

	if idErr != nil {
		return entity.Spectrum{}, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}

	spectrumData := entity.Spectrum{
		Id:               spectrumId,
		Name:             name,
		ShortDescription: shortDescription,
		Description:      description,
		Status:           enums.Pending,
		ImageUrls:        imageUrls,
		LogoUrl:          logoUrl,
		VideoUrl:         videoUrl,
		TotalEvents:      0,
	}

	spectrum, creationErr := s.repo.CreateSpectrum(spectrumData)

	if creationErr != nil {
		return entity.Spectrum{}, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}

	return spectrum, nil
}
