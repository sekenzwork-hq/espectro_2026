package repository

import (
	"espectro/entity"

	"gorm.io/gorm"
)

type AdminPostgresRepo struct {
	db *gorm.DB
}

func NewAdminPostgresRepo(db *gorm.DB) AdminPostgresRepo {
	return AdminPostgresRepo{
		db: db,
	}
}

func (a AdminPostgresRepo) RetrieveAdminCredByEmail(email string) (entity.AdminDBLoginCredentials, error) {

	cred := &entity.AdminDBLoginCredentials{
		Email: email,
	}
	err := a.db.Table("admins").Where("email=? AND deleted_at IS NULL", email).Select("password", "id").First(&cred).Error
	return *cred, err
}
