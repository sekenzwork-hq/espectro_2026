package repository

type EventSponsors interface {
	AddSponsor(sponsorId string, eventIds []string) error
}
