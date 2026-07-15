package repository

import (
	"espectro/entity"
)

type PartnerRepo interface {
	AddPartner(partner entity.PartnerEntity) (entity.PartnerEntity, error)
}
