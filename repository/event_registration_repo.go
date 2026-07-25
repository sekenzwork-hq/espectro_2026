package repository

import "espectro/entity"

type EventRegistrationRepo interface {
	Register(registrationDetails entity.EventRegistrationEntity) (entity.EventRegistrationEntity, error)
	RegisterExists(eventId string, userId string) (bool, error)
}
