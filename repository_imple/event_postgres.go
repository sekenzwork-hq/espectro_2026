package repositoryimple

import (
	"espectro/entity"

	"gorm.io/gorm"
)

type EventPostgresRepo struct {
	db *gorm.DB
}

func (e EventPostgresRepo) CreateEvent(event entity.EventEntity) (entity.EventEntity, error) {
	err := e.db.Create(&event).Error
	return event, err
}
