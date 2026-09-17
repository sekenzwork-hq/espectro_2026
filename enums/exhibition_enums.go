package enums

type ExhibitionStatus string

const (
	ApprovedExhibition ExhibitionStatus = "approved"
	PendingExhibition  ExhibitionStatus = "pending"
	RejectedExhibition ExhibitionStatus = "rejected"
)

func (e ExhibitionStatus) IsValid() bool {

	switch e {
	case ApprovedExhibition, PendingExhibition, RejectedExhibition:
		return true
	default:
		return false
	}
}
