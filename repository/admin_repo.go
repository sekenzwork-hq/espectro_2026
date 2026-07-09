package repository

import (
	"espectro/entity"
	"espectro/enums"

	"github.com/google/uuid"
)

type AdminRepo interface {
	RetrieveAdminCredByEmail(email string) (entity.AdminDBLoginCredentials, error)
	CreateNewAdmin(admin entity.AdminEntity) (uuid.UUID, error)
	RetrieveAdminRoleByID(adminId string) (string, error)
	DeleteMemberOrVolunteer(adminId string) error
	CheckAdminExists(adminId string) (bool, error)
	UpdateCurrentAdmin(adminId string, newEmail string, newFullname string, updateMode enums.AdminUpdateMode) error
}
