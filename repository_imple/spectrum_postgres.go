package repositoryimple

import (
	"espectro/entity"

	"gorm.io/gorm"
)

type SpectrumPostgresRepo struct {
	db *gorm.DB
}

func NewSpectrumPostgresRepo(db *gorm.DB) SpectrumPostgresRepo {
	return SpectrumPostgresRepo{db: db}
}

func (s SpectrumPostgresRepo) CreateSpectrum(newSpectrum entity.Spectrum) (entity.Spectrum, error) {
	out := s.db.Table("spectrums").Create(&newSpectrum)
	return newSpectrum, out.Error
}
