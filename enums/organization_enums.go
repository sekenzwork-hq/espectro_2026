package enums

type OrganizationStatus string

const (
	ApprovedOrganization  OrganizationStatus = "approved"
	PendingOrganization   OrganizationStatus = "pending"
	RejectedOrganization  OrganizationStatus = "rejected"
	ForBiddenOrganization OrganizationStatus = "forbidden"
)

func (o OrganizationStatus) IsValid() bool {
	switch o {
	case ApprovedOrganization, PendingOrganization, RejectedOrganization, ForBiddenOrganization:
		return true
	default:
		return true
	}
}
