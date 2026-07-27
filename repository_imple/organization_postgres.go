package repositoryimple

import (
	"espectro/entity"

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
