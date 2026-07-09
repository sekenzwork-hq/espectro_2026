package enums

type AdminUpdateMode string

const (
	EmailOnly        AdminUpdateMode = "EMAIL_ONLY"
	FullnameOnly     AdminUpdateMode = "FULLNAME_ONLY"
	EmailAndFullname AdminUpdateMode = "EMAIL_AND_FULLNAME"
)

func (a AdminUpdateMode) IsValid() bool {

	switch a {
	case EmailAndFullname, EmailOnly, FullnameOnly:
		return true
	default:
		return false
	}
}
