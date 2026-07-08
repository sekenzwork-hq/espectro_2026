package repository

import "espectro/entity"

type AdminRepo interface {
	RetrieveAdminCredByEmail(email string) (entity.AdminDBLoginCredentials, error)
}
