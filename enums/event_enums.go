package enums

type EventStatus string

const (
	PendingEvent   EventStatus = "pending"
	OngoingEvent   EventStatus = "on_going"
	CancelledEvent EventStatus = "cancelled"
	ScheduledEvent EventStatus = "scheduled"
)

func (e EventStatus) IsValid() bool {

	switch e {
	case PendingEvent, OngoingEvent, CancelledEvent, ScheduledEvent:
		return true
	default:
		return false
	}
}

type EventMode string

const (
	OnlineEvent  EventMode = "online"
	StadiumEvent EventMode = "stadium"
	HallEvent    EventMode = "hall"
)

func (e EventMode) IsValid() bool {

	switch e {
	case OnlineEvent, StadiumEvent, HallEvent:
		return true
	default:
		return false
	}
}

type EventType string

const (
	SpeechEvent      EventType = "event"
	ConcertEvent     EventType = "concert"
	CompetitionEvent EventType = "competition"
)

func (e EventType) IsValid() bool {

	switch e {
	case SpeechEvent, ConcertEvent, CompetitionEvent:
		return true
	default:
		return false
	}
}
