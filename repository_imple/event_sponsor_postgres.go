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
