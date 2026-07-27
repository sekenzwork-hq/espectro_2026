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
			participant_limit=COALESCE(?,participant_limit),
			start_date=COALESCE(?,start_date),
			end_date=COALESCE(?,end_date),
			event_mode=COALESCE(?,event_mode),
			event_type=COALESCE(?,event_type),
			is_featured=COALESCE(?,is_featured),
			contact_email=COALESCE(?,contact_email),
			venue_id=COALESCE(?,venue_id) WHERE id=? AND deleted_at IS NULL

			RETURNING 
		    id,name,description,spectrum_id,status,participant_limit,start_date,end_date,event_mode,event_type,is_featured,contact_email,total_registrations,venue_id,created_at
			`, newEvent.Name,
			newEvent.Description,
			newEvent.SpectrumId,
			newEvent.Status,
			newEvent.ParticipantLimit,
			newEvent.StartDate,
			newEvent.EndDate,
			newEvent.EventMode,
			newEvent.EventType,
			newEvent.IsFeatured,
			newEvent.ContactEmail,
			newEvent.VenueId,
			eventId,
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
		Select("id,name,description,status,event_mode,event_type,participant_limit,start_date,end_date,spectrum_id,total_registrations,venue_id,contact_email,created_at").
		Where("deleted_at IS NULL").
		Limit(limit).
		Offset(offset).
		Scan(&events).Error

	return events, err
}

func (e EventPostgresRepo) CheckMultipleEventsExist(eventIds []string) (bool, error) {

	var count int64
	err := e.db.
		Raw(`SELECT COUNT(*) FROM events WHERE id=ANY(?)`, eventIds).
		Scan(&count).Error

	return int(count) == len(eventIds), err
}

func (e EventPostgresRepo) RetrieveSpectrumId(eventId string) (string, error) {

	var id string
	err := e.db.
		Table("events").
		Select("spectrum_id").
		Where("id=? AND deleted_at IS NULL").
		Scan(&id).Error

	if err != nil {
		return "", err
	} else if id == "" {
		return "", gorm.ErrRecordNotFound
	}

	return id, nil
}
func (e EventPostgresRepo) DeleteEventBySpectrumId(spectrumId string) error {

	out := e.db.
		Table("events").
		Where("spectrum_id=? AND deleted_at IS NULL", spectrumId).
		UpdateColumn("deleted_at", time.Now().UTC())

	if out.Error != nil {
		return out.Error
	} else if out.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (e EventPostgresRepo) RetrieveEventsBySpectrumId(spectrumId string, limit int, offset int) ([]entity.EventEntity, error) {
	var events []entity.EventEntity

	err := e.db.
		Table("events").
		Select("id,name,description,spectrum_id,status,start_date,end_date,participant_limit,event_mode,event_type,is_featured,contact_email,venue_id,created_at").
		Where("spectrum_id=? AND deleted_at IS NULL", spectrumId).
		Offset(offset).
		Limit(limit).
		Scan(&events).Error

	return events, err
}

func (e EventPostgresRepo) EventExists(eventId string) (bool, error) {

	var exists bool

	err := e.db.Raw(
		`SELECT EXISTS (SELECT 1 FROM events WHERE id=? AND deleted_at IS NULL)`,
		eventId,
	).Scan(&exists).Error

	return exists, err
}

func (e EventPostgresRepo) IncrementTotalRegistrationBy1(id string) error {

	out := e.db.
		Table("events").
		Where("id=? AND deleted_at IS NULL", id).
		Update("total_registrations", gorm.Expr("total_registrations+1"))

	if out.Error != nil {
		return out.Error
	} else if out.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (e EventPostgresRepo) DecrementTotalRegistrationBy1(id string) error {

	out := e.db.Table("events").Where("id=? AND deleted_at IS NULL", id).Update("total_registrations", gorm.Expr(`
	CASE 
		WHEN total_registrations != 0 THEN total_registrations-1
		ELSE total_registrations
	END
	`))

	if out.Error != nil {
		return out.Error
	} else if out.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
