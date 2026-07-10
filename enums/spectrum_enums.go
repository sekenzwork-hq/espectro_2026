package enums

type SpectrumStatus string

const (
	Pending   SpectrumStatus = "pending"
	OnGoing   SpectrumStatus = "on_going"
	Scheduled SpectrumStatus = "scheduled"
)

func (s SpectrumStatus) IsValid() bool {

	switch s {
	case Pending, OnGoing, Scheduled:
		return true
	default:
		return false
	}
}
