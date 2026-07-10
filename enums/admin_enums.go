package enums

type AdminUpdateMode string

const (
	EmailOnly        AdminUpdateMode = "email"
	FullnameOnly     AdminUpdateMode = "fullname_only"
	EmailAndFullname AdminUpdateMode = "email_and_fullname"
)

func (a AdminUpdateMode) IsValid() bool {

	switch a {
	case EmailAndFullname, EmailOnly, FullnameOnly:
		return true
	default:
		return false
	}
}

type AdminRole string

const (
	Leader    AdminRole = "leader"
	Member    AdminRole = "member"
	Volunteer AdminRole = "volunteer"
)

func (a AdminRole) IsValid() bool {
	switch a {
	case Leader, Member, Volunteer:
		return true
	default:
		return false
	}
}

func (a AdminRole) ParseRole(strRole string) (AdminRole, bool) {

	switch strRole {
	case "leader":
		return Leader, true
	case "member":
		return Member, true
	case "volunteer":
		return Volunteer, true
	default:
		return "", false
	}
}

type AdminMiddlewareType string

const (
	LeaderMiddleware          AdminMiddlewareType = "leader_middleware"
	AllAdminMiddleware        AdminMiddlewareType = "all_admin_middleware"
	LeaderAndMemberMiddleware AdminMiddlewareType = "leader_and_member"
)

func (a AdminMiddlewareType) IsValid() bool {

	switch a {
	case LeaderMiddleware, AllAdminMiddleware, LeaderAndMemberMiddleware:
		return true
	default:
		return false
	}
}
