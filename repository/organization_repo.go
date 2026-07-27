package repository

import "espectro/entity"

type OrganizationRepo interface {
	CreateOrganization(organization entity.OrganizationEntity) (entity.OrganizationEntity, error)
}
