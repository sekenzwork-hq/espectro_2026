package repositoryimple

import (
	"espectro/entity"
	"espectro/enums"

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
		`SELECT EXISTS (SELECT 1 FROM event_registrations WHERE event_id=? AND user_id=? AND deleted_at IS NULL)`,
		eventId, userId,
	).Scan(&exists).Error
	return exists, err
}
func (e EventRegistrationPostgresRepo) ChangeRegistrationStatus(registrationId string, newStatus enums.EventRegStatus) (entity.EventRegistrationEntity, error) {
	var eventRegistration entity.EventRegistrationEntity
	out := e.db.Raw(
		`
		UPDATE event_registrations SET 
		status=?,
		check_in=now()
		WHERE id=? AND deleted_at IS NULl
		RETURNING id,user_id,event_id,status,check_in,check_out,created_at
		`,
		newStatus, registrationId,
	).Scan(&eventRegistration)

	if out.Error != nil {
		return entity.EventRegistrationEntity{}, out.Error
	} else if out.RowsAffected == 0 {
		return entity.EventRegistrationEntity{}, gorm.ErrRecordNotFound
	}

	return eventRegistration, nil
}
