package repositoryimple

import (
	"espectro/entity"
	"time"

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

func (p PartnerPostgresRepo) UpdatePartner(partner entity.PartnerUpdateEntity) error {

	data := map[string]any{}

	if partner.LogoUrl != nil {
		data["logo_url"] = *partner.LogoUrl
	}
	if partner.Name != nil {
		data["name"] = *partner.Name
	}

	out := p.db.
		Table("partners").
		Where("id=? AND deleted_at IS NULL", partner.Id).
		Updates(data)

	if out.Error != nil {
		return out.Error
	} else if out.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (p PartnerPostgresRepo) DeletePartner(partnerId string) error {

	out := p.db.
		Table("partners").
		Where("id=? AND deleted_at IS NULL", partnerId).
		Update("deleted_at", time.Now().UTC())

	if out.Error != nil {
		return out.Error
	} else if out.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (p PartnerPostgresRepo) RetrievePartner(limit int, offset int) ([]entity.PartnerEntity, error) {

	var partners []entity.PartnerEntity

	err := p.db.
		Table("partners").
		Select("id,name,logo_url,created_at").
		Where("deleted_at IS NULL").
		Limit(limit).
		Offset(offset).
		Scan(&partners).Error

	return partners, err

}
