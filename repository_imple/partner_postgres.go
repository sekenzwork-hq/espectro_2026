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

func (p PartnerPostgresRepo) CreatePartner(partner entity.PartnerEntity) (entity.PartnerEntity, error) {
	err := p.db.
		Table("partners").
		Create(&partner).Error

	return partner, err
}

func (p PartnerPostgresRepo) UpdatePartner(partnerId string, name *string, logoUrl *string) (entity.PartnerEntity, error) {

	var newPartner entity.PartnerEntity

	out := p.db.
		Raw(`UPDATE partners 
		SET name = COALESCE(?,name), 
		    logo_url = COALESCE(?,logo_url)
		WHERE id=? AND deleted_at IS NULL RETURNING id,name,logo_url,created_at`, name, logoUrl, partnerId).
		Scan(&newPartner)

	if out.Error != nil {
		return entity.PartnerEntity{}, out.Error
	} else if out.RowsAffected == 0 {
		return entity.PartnerEntity{}, gorm.ErrRecordNotFound
	}

	return newPartner, nil
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
