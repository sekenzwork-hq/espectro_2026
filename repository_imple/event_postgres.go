package repositoryimple

import (
	"espectro/entity"
	"time"

	"gorm.io/gorm"
)

type EventPostgresRepo struct {
	db *gorm.DB
}

func NewEventPostgresRepo(db *gorm.DB) EventPostgresRepo {
	return EventPostgresRepo{db: db}
}
func (e EventPostgresRepo) CreateEvent(event entity.EventEntity) (entity.EventEntity, error) {
	err := e.db.
		Table("events").
		Create(&event).Error
	return event, err
}

func (e EventPostgresRepo) UpdateEvent(eventId string, newEvent entity.EventUpdateEntity) error {
	out := e.db.
		Table("events").
		Where("id=? AND deleted_at IS NULL", eventId).
		Updates(&newEvent)

	if out.Error != nil {
		return out.Error
	} else if out.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	} else {
		return nil
	}
}

func (e EventPostgresRepo) DeleteEvent(eventId string) error {

	out := e.db.
		Table("events").
		Where("id=? AND deleted_at IS NULL", eventId).UpdateColumn("deleted_at", time.Now().UTC())

	if out.Error != nil {
		return out.Error
	} else if out.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (e EventPostgresRepo) RetrieveEvents(limit int, offset int) ([]entity.EventEntity, error) {

	var events []entity.EventEntity

	err := e.db.
		Table("events").
		Select("id,name,description,status,event_mode,event_type,participant_limit,start_date,end_date,spectrum_id,venue_id,contact_email,created_at").
		Where("deleted_at IS NULL").
		Limit(limit).
		Offset(offset).
		Scan(&events).Error

	return events, err
}
