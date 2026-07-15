package repositoryimple

import (
	"espectro/entity"

	"gorm.io/gorm"
)

type PartnerPostgresRepo struct {
	db *gorm.DB
}

func NewPartnerPostgres(db *gorm.DB) PartnerPostgresRepo {
	return PartnerPostgresRepo{db: db}
}

func (p PartnerPostgresRepo) AddPartner(partner entity.PartnerEntity) (entity.PartnerEntity, error) {

	err := p.db.
		Table("partners").
		Create(&partner).Error

	return partner, err
}
