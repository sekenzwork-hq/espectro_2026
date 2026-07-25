package enums

type StaffRole string

const (
	Organizer StaffRole = "organizer"
	Manager   StaffRole = "manager"
	Supplier  StaffRole = "supplier"
	All       StaffRole = "all"
)

func (s StaffRole) IsValid() bool {
	switch s {
	case Organizer, Manager, Supplier, All:
		return true
	default:
		return false
	}
}
