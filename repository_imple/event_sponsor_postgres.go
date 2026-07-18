package repositoryimple

import (
	"espectro/entity"

	"gorm.io/gorm"
)

type EventSponsorPostgresRepo struct {
	db *gorm.DB
}

func NewEventSponsorPostgresRepo(db *gorm.DB) EventPostgresRepo {
	return EventPostgresRepo{db: db}
}

func (e EventPostgresRepo) AddSponsor(sponsorId string, eventIds []string) error {

	records := make([]entity.EventSponsorEntity, len(eventIds))
	for i := range eventIds {
		eventId := eventIds[i]
		records = append(records, entity.EventSponsorEntity{
			SponsorId: sponsorId,
			EventId:   eventId,
		})
	}
	err := e.db.
		Table("event_sponsors").
		Create(&records).Error

	return err
}

func (e EventPostgresRepo) RetrieveSponsorsBasedOnEvent(eventId string, limit int, offset int) ([]entity.SponsorEntity, error) {

	var sponsors []entity.SponsorEntity

	err := e.db.Raw(
		`
		SELECT s.id,s.name,s.amount,s.profile_or_org_url,s.sponsored_type,s.created_at FROM sponsors s
		INNER JOIN event_sponsors es ON s.id=es.sponsor_id WHERE event_id=?
		GROUP BY s.id OFFSET ? LIMIT ?
		`,
		eventId, offset, limit,
	).
		Offset(offset).
		Limit(limit).
		Scan(&sponsors).Error

	return sponsors, err
}
