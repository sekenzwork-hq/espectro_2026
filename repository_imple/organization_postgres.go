package repositoryimple

import (
	"espectro/entity"
	"time"

	"gorm.io/gorm"
)

type OrganizationPostgresRepo struct {
	db *gorm.DB
}

func NewOrganizationPostgresRepo(db *gorm.DB) OrganizationPostgresRepo {
	return OrganizationPostgresRepo{db: db}
}

func (o OrganizationPostgresRepo) CreateOrganization(organization entity.OrganizationEntity) (entity.OrganizationEntity, error) {

	err := o.db.
		Table("organizations").
		Create(&organization).Error

	return organization, err
}

func (o OrganizationPostgresRepo) UpdateOrganization(newOrganization entity.OrganizationUpdateEntity) (entity.OrganizationEntity, error) {
	var org entity.OrganizationEntity

	out := o.db.
		Raw(
			`
			UPDATE organizations SET
			fullname=COALESCE(?,fullname),
			email=COALESCE(?,email),
			logo_url=COALESCE(?,logo_url),
			status=COALESCE(?,status),
			phone=COALESCE(?,phone),
			website_url=COALESCE(?,website_url),
			industry=COALESCE(?,industry),
			headquarters=COALESCE(?,headquarters)
            
			WHERE id=? AND deleted_at IS NULL
			RETURNING id,fullname,email,logo_url,status,phone,website_url,industry,headquarters

			`,
			newOrganization.Fullname,
			newOrganization.Email,
			newOrganization.LogoUrl,
			newOrganization.Status,
			newOrganization.PhoneNumber,
			newOrganization.WebsiteUrl,
			newOrganization.Industry,
			newOrganization.HeadQuarters,
			newOrganization.Id,
		).Scan(&org)

	if out.Error != nil {
		return entity.OrganizationEntity{}, out.Error
	} else if out.RowsAffected == 0 {
		return entity.OrganizationEntity{}, gorm.ErrRecordNotFound
	}

	return org, nil
}

func (o OrganizationPostgresRepo) DeleteOrganization(id string) error {

	out := o.db.
		Table("organizations").
		Where("id=? AND deleted_at IS NULL", id).
		Update("deleted_at", time.Now().UTC())

	if out.Error != nil {
		return out.Error
	} else if out.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (o OrganizationPostgresRepo) RetrieveOrganizations(limit int, offset int) ([]entity.OrganizationEntity, error) {

	var organizations []entity.OrganizationEntity

	err := o.db.
		Table("organizations").
		Where("deleted_at IS NULL").
		Offset(offset).Limit(limit).
		Scan(&organizations).Error

	return organizations, err
}
