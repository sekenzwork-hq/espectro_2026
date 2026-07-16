package repository

import "espectro/entity"

type EventRepo interface {
	CreateEvent(event entity.EventEntity) (entity.EventEntity, error)
	UpdateEvent(eventId string, newEvent entity.EventUpdateEntity) (entity.EventEntity, error)
	DeleteEvent(eventId string) error
	RetrieveEvents(limit int, offset int) ([]entity.EventEntity, error)
}
