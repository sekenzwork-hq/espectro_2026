package enums

type SpectrumStatus string

const (
	PendingSpectrum   SpectrumStatus = "pending"
	OnGoingSpectrum   SpectrumStatus = "on_going"
	ScheduledSpectrum SpectrumStatus = "scheduled"
)

func (s SpectrumStatus) IsValid() bool {

	switch s {
	case PendingSpectrum, OnGoingSpectrum, ScheduledSpectrum:
		return true
	default:
		return false
	}
}
