package repository

import (
	"espectro/entity"
	"espectro/enums"
)

type SpectrumRepo interface {
	CreateSpectrum(newSpectrum entity.SpectrumEntity) (entity.SpectrumEntity, error)
	UpdateSpectrum(
		spectrumId string,
		name *string,
		shortDescription *string,
		description *string,
		status *enums.SpectrumStatus,
		logoUrl *string,
		videoUrl *string,
		imageUrls []string,
	) error
	DeleteSpectrum(spectrumId string) error
	RetrieveSpectrums(offset int, limit int) ([]entity.SpectrumEntity, error)
	CheckSpectrumExists(spectrumId string) (bool, error)
}
