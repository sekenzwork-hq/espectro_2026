package enums

type CompetitionStatusEnum string

const (
	PendingCompetition   CompetitionStatusEnum = "pending"
	OnGoingCompetition   CompetitionStatusEnum = "on_going"
	FinishedCompetition  CompetitionStatusEnum = "finished"
	ScheduledCompetition CompetitionStatusEnum = "scheduled"
	CancelledCompetition CompetitionStatusEnum = "cancelled"
)

func (c CompetitionStatusEnum) IsValid() bool {

	switch c {

	case PendingCompetition, OnGoingCompetition, FinishedCompetition, ScheduledCompetition, CancelledCompetition:
		return true
	default:
		return false
	}
}
