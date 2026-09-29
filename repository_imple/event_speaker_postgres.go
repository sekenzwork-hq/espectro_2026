package repositoryimple

import (
	"gorm.io/gorm"
)

type EventSpeakerPostgresRepo struct {
	db *gorm.DB
}

func NewEventSpeakerPostgresRepo(db *gorm.DB) EventSpeakerPostgresRepo {
	return EventSpeakerPostgresRepo{db: db}
}

func (e EventSpeakerPostgresRepo) AddSpeakerToEvent(speakerId string, eventId string) error {

	err := e.db.Exec(
		`
		INSERT INTO event_speakers(event_id,speaker_id) VALUES(?,?)
		`,
		eventId,
		speakerId,
	).Error

	return err
}
