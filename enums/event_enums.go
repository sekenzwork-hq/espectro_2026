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
