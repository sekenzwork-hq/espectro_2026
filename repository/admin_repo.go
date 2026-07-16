package repository

import (
	"espectro/entity"
	"espectro/enums"

	"github.com/google/uuid"
)

type AdminRepo interface {
	RetrieveAdminCredByEmail(email string) (entity.AdminDBLoginCredentials, error)
	CreateNewAdmin(admin entity.AdminCreateEntity) (uuid.UUID, error)
	RetrieveAdminRoleByID(adminId string) (enums.AdminRole, error)
	DeleteMemberOrVolunteer(adminId string) error
	CheckAdminExists(adminId string) (bool, error)
	UpdateCurrentAdmin(adminId string, admin entity.AdminUpdateEntity) (entity.AdminEntity, error)
	UpdateAdminRole(adminId string, newRole enums.AdminRole) error
}
