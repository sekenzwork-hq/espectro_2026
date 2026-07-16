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

func (e EventPostgresRepo) UpdateEvent(eventId string, newEvent entity.EventUpdateEntity) (entity.EventEntity, error) {

	var event entity.EventEntity

	out := e.db.
		Raw(
			`UPDATE events SET 
			name=COALESCE(?,name),
			description=COALESCE(?,description),
			spectrum_id=COALESCE(?,spectrum_id),
			status=COALESCE(?,status),
			limit=COALESCE(?,limit),
			start_date=COALESCE(?,start_date),
			end_date=COALESCE(?,end_date),
			event_mode=COALESCE(?,event_mode),
			event_type=COALESCE(?,event_type),
			is_featured=COALESCE(?,is_featured),
			contact_email=COALESCE(?,contact_email),
			venue_id=COALESCE(?,venue_id) 

			RETURNING 
		    name,description,spectrum_id,status,limit,start_date,end_date,event_mode,event_type,is_featured,contact_email,venue_id
			`,
		).
		Scan(&event)

	if out.Error != nil {
		return event, out.Error
	} else if out.RowsAffected == 0 {
		return event, gorm.ErrRecordNotFound
	}
	return event, nil

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
