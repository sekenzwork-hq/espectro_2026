package repository

import (
	"espectro/entity"

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
	cred := entity.AdminDBLoginCredentials{
		Email: email,
	}
	err := a.db.
		Table("admins").
		Where("email=? AND deleted_at IS NULL", email).
		Select("password", "id").
		First(&cred).Error

	return cred, err
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

func (a AdminPostgresRepo) DeleteMemberOrVolunteer(adminId string) error {
	out := a.db.
		Table("admins").
		Where("id=? AND (admin_role = 'member' OR admin_role = 'volunteer')", adminId).
		Delete(nil)

	if out.Error != nil {
		return out.Error
	} else if out.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (a AdminPostgresRepo) CheckAdminExists(adminId string) (bool, error) {

	var exists bool
	err := a.db.
		Raw("SELECT EXISTS(SELECT 1 FROM admins WHERE id=?)", adminId).
		Scan(&exists).Error

	return exists, err

}
