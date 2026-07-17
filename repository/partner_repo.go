package repository

import (
	"espectro/entity"
)

type PartnerRepo interface {
	CreatePartner(partner entity.PartnerEntity) (entity.PartnerEntity, error)
	UpdatePartner(partnerId string, name *string, logoUrl *string) (entity.PartnerEntity, error)
	DeletePartner(partnerId string) error
	RetrievePartner(limit int, offset int) ([]entity.PartnerEntity, error)
}
