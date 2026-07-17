package repository

import "espectro/entity"

type SponsorRepo interface {
	CreateSponsor(sponsor entity.SponsorEntity) (entity.SponsorEntity, error)
}
