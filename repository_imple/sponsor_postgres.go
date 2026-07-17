package repositoryimple

import (
	"espectro/entity"

	"gorm.io/gorm"
)

type SponsorPostgresRepo struct {
	db *gorm.DB
}

func NewSponsorPostgresRepo(db *gorm.DB) SponsorPostgresRepo {
	return SponsorPostgresRepo{db: db}
}

func (s SponsorPostgresRepo) CreateSponsor(sponsor entity.SponsorEntity) (entity.SponsorEntity, error) {

	err := s.db.
		Table("sponsors").
		Create(&sponsor).Error

	return sponsor, err
}
