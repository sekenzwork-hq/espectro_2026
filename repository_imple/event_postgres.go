package repositoryimple

import (
	"espectro/entity"

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
