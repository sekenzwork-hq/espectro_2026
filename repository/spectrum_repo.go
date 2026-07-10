package repository

import "espectro/entity"

type SpectrumRepo interface {
	CreateSpectrum(newSpectrum entity.Spectrum) (entity.Spectrum, error)
}
