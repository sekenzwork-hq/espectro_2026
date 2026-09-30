package repositoryimple

import (
	"time"

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

func (e EventSpeakerPostgresRepo) SpeakerIsAlreadyAddedToEvent(speakerId string, eventId string) (bool, error) {

	var added bool

	err := e.db.Raw(
		`SELECT EXISTS (SELECT 1 FROM event_speakers WHERE event_id=? AND speaker_id=? AND deleted_at IS NULL)`,
		eventId,
		speakerId,
	).Scan(&added).Error

	return added, err
}

func (e EventSpeakerPostgresRepo) RemoveSpeakerFromEvent(speakerId string, eventId string) error {

	out := e.db.
		Table("event_speakers").
		Where("event_id=? AND speaker_id=? AND deleted_at IS NULL", eventId, speakerId).
		UpdateColumn("deleted_at", time.Now().UTC())

	if out.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return out.Error
}
