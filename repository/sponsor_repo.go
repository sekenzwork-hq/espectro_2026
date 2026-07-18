package repository

import "espectro/entity"

type SponsorRepo interface {
	CreateSponsor(sponsor entity.SponsorEntity) (entity.SponsorEntity, error)
	UpdateSponsor(sponsor entity.SponsorUpdateEntity) (entity.SponsorEntity, error)
	DeleteSponsor(sponsorId string) error
	RetrieveSponsors(limit int, offset int) ([]entity.SponsorEntity, error)
}
