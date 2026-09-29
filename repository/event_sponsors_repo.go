package repository

import "espectro/entity"

type EventSponsors interface {
	AddSponsor(sponsorId string, eventIds []string) error
	RetrieveSponsorsBasedOnEvent(eventId string, offset int) ([]entity.SponsorEntity, error)
}
