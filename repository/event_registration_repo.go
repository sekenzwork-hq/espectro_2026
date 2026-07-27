package repository

import (
	"espectro/entity"
	"espectro/enums"
)

type EventRegistrationRepo interface {
	Register(registrationDetails entity.EventRegistrationEntity) (entity.EventRegistrationEntity, error)
	RegisterExists(eventId string, userId string) (bool, error)
	ChangeRegistrationStatus(registrationId string, newStatus enums.EventRegStatus) (entity.EventRegistrationEntity, error)
	RetrieveEventIdUsingRegistrationId(id string) (string, error)
	WithdrawRegistration(id string) error
}
