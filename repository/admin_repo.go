package repository

import "espectro/entity"

type AdminRepo interface {
	RetrieveAdminCredByEmail(email string) (entity.AdminDBLoginCredentials, error)
	CreateNewAdmin(admin entity.AdminCreateEntity) error
	RetrieveAdminRoleByID(adminId string) (string, error)
}
