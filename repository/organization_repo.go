package repository

import "espectro/entity"

type OrganizationRepo interface {
	CreateOrganization(organization entity.OrganizationEntity) (entity.OrganizationEntity, error)
	UpdateOrganization(newOrganization entity.OrganizationUpdateEntity) (entity.OrganizationEntity, error)
	DeleteOrganization(id string) error
	RetrieveOrganizations(limit int, offset int) ([]entity.OrganizationEntity, error)
	OrganizationExists(id string) (bool, error)
}
