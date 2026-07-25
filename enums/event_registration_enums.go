package enums

type EventRegStatus string

const (
	VerificationPending EventRegStatus = "verification_pending"
	Approved            EventRegStatus = "Approved"
	Rejected            EventRegStatus = "rejected"
)

func (e EventRegStatus) IsValid() bool {

	switch e {
	case VerificationPending, Approved, Rejected:
		return true
	default:
		return false
	}
}
