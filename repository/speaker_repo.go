package repository

import "espectro/entity"

type SpeakerRepo interface {
	CreateSpeaker(newSpeaker entity.SpeakerDBCreateEntity) (entity.SpeakerDBRetrieveEntity, error)
	UpdateSpeaker(speakerId string, newSpeaker entity.SpeakerDBUpdateEntity) (entity.SpeakerDBRetrieveEntity, error)
	DeleteSpeaker(speakerId string) error
	SpeakerExists(speakerId string) (bool, error)
	RetrieveSpeaker(eventId string, offset int) ([]entity.SpeakerDBRetrieveEntity, error)
}
