package repositoryimple

import (
	"espectro/entity"

	"gorm.io/gorm"
)

type EventRegistrationPostgresRepo struct {
	db *gorm.DB
}

func NewEventRegistrationPostgresRepo(db *gorm.DB) EventRegistrationPostgresRepo {
	return EventRegistrationPostgresRepo{db: db}
}

func (e EventRegistrationPostgresRepo) Register(registrationDetails entity.EventRegistrationEntity) (entity.EventRegistrationEntity, error) {
	err := e.db.
		Table("event_registrations").
		Create(&registrationDetails).Error

	return registrationDetails, err
}

func (e EventRegistrationPostgresRepo) RegisterExists(eventId string, userId string) (bool, error) {
	var exists bool
	err := e.db.Raw(
		`SELECT EXISTS (SELECT 1 FROM event_registration WHERE event_id=? AND user_id=? AND deleted_at IS NULL)`,
		eventId, userId,
	).Scan(&exists).Error
	return exists, err
}
