package repository

import "espectro/entity"

type EventRepo interface {
	CreateEvent(event entity.EventEntity) (entity.EventEntity, error)
}
