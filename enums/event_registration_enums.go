package enums

type EventRegStatus string

const (
	VerificationPendingEventRegistration EventRegStatus = "verification_pending"
	ApprovedEventRegistration            EventRegStatus = "approved"
	RejectedEventRegistration            EventRegStatus = "rejected"
)

func (e EventRegStatus) IsValid() bool {

	switch e {
	case VerificationPendingEventRegistration, ApprovedEventRegistration, RejectedEventRegistration:
		return true
	default:
		return false
	}
}
