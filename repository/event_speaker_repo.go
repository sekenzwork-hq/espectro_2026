package repository

type EventSpeakerRepo interface {
	SpeakerIsAlreadyAddedToEvent(speakerId string, eventId string) (bool, error)
	AddSpeakerToEvent(speakerId string, eventId string) error
	RemoveSpeakerFromEvent(speakerId string, eventId string) error
}
