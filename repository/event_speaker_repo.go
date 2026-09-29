package repository

type EventSpeakerRepo interface {
	AddSpeakerToEvent(speakerId string, eventId string) error
}
