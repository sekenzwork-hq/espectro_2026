package usecases

import (
	customerrors "espectro/custom_errors"
	"espectro/entity"
	"espectro/pkg"
	"espectro/repository"

	"github.com/lib/pq"
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
	imageUrl []string,
	videoUrl *string,
	logoUrl *string,

) (entity.Spectrum, error) {

	nameErr := pkg.ValidateSpectrumOrEventName(name)

	if nameErr != nil {
		return entity.Spectrum{}, nameErr
	}

	shortDescriptionErr := pkg.ValidateSpectrumShortDescription(shortDescription)

	if shortDescriptionErr != nil {
		return entity.Spectrum{}, nameErr
	}

	descriptionErr := pkg.ValidateSpectrumOrEventDescription(description)

	if descriptionErr != nil {
		return entity.Spectrum{}, descriptionErr
	}

	if len(imageUrl) > 10 {
		return entity.Spectrum{}, &customerrors.SpaceError{OrgError: "Maximum number of images is 10"}
	}

	spectrum, creationErr := s.repo.CreateSpectrum(entity.Spectrum{
		Name:             name,
		ShortDescription: shortDescription,
		Description:      description,
		LogoUrl:          nil,
		ImageUrls:        pq.StringArray{},
		VideoUrl:         nil,
	})

	if creationErr != nil {
		return entity.Spectrum{}, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}

	return spectrum, nil
}
