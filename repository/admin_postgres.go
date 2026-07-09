package repository

import (
	"espectro/entity"
	"time"

	"github.com/google/uuid"
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
	err := a.db.
		Table("admins").
		Where("email=? AND deleted_at IS NULL", email).
		Select("password", "id").
		First(&cred).Error
	return *cred, err
}

func (a AdminPostgresRepo) CreateNewAdmin(admin entity.AdminEntity) (uuid.UUID, error) {
	return admin.Id, a.db.
		Table("admins").
		Create(&admin).Error
}

func (a AdminPostgresRepo) RetrieveAdminRoleByID(adminId string) (string, error) {
	var role string
	err := a.db.
		Table("admins").
		Where("id=?", adminId).
		Select("admin_role").
		Pluck("admin_role", &role).Error
	return role, err
}

func (a AdminPostgresRepo) SoftDeleteMemberOrVolunteer(adminId string) error {
	return a.db.
		Table("admins").
		Where("id=?", adminId).
		Update("deleted_at", time.Now().UTC()).Error
}
