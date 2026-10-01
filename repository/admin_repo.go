package repository

import (
	"espectro/entity"
	"espectro/enums"
)

type AdminRepo interface {
	RetrieveAdminCredByEmail(email string) (entity.AdminDBLoginCredentials, error)
	CreateNewAdmin(admin entity.AdminCreateEntity) (entity.AdminEntity, error)
	RetrieveAdminRoleByID(adminId string) (enums.AdminRole, error)
	DeleteMemberOrVolunteer(adminId string) error
	AdminExists(adminId string) (bool, error)
	AdminEmailExists(email string) (bool, error)
	UpdateCurrentAdmin(adminId string, admin entity.AdminUpdateEntity) (entity.AdminEntity, error)
	UpdateAdminRole(adminId string, newRole enums.AdminRole) error
}
