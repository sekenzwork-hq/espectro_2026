package repositoryimple

import (
	"espectro/entity"
	"espectro/enums"

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

func (a AdminPostgresRepo) CreateNewAdmin(admin entity.AdminCreateEntity) (uuid.UUID, error) {
	return admin.Id, a.db.
		Table("admins").
		Create(&admin).Error
}

func (a AdminPostgresRepo) RetrieveAdminRoleByID(adminId string) (enums.AdminRole, error) {
	var role string
	err := a.db.
		Table("admins").
		Where("id=? AND deleted_at IS NULL", adminId).
		Select("admin_role").
		Pluck("admin_role", &role).Error

	if err != nil {
		return "", err
	} else if role == "" {
		return "", gorm.ErrRecordNotFound
	} else {
		return enums.AdminRole(role), nil
	}

}

func (a AdminPostgresRepo) DeleteMemberOrVolunteer(adminId string) error {
	out := a.db.
		Table("admins").
		Where("id=? AND deleted_at IS NULL AND (admin_role = 'member' OR admin_role = 'volunteer')", adminId).
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
		Raw("SELECT EXISTS(SELECT 1 FROM admins WHERE id=? AND deleted_at IS NULL)", adminId).
		Scan(&exists).Error

	return exists, err

}

func (a AdminPostgresRepo) UpdateCurrentAdmin(adminId string, admin entity.AdminUpdateEntity) (entity.AdminEntity, error) {

	var newAdmin entity.AdminEntity
	out := a.db.
		Raw(
			`
			UPDATE admins SET fullname=COALESCE(?,fullname),email=COALESCE(?,email)
			WHERE id=? AND deleted_at IS NULL 
			RETURNING id,fullname,email,admin_role,created_at
			`, admin.Fullname, admin.Email, adminId,
		).Scan(&newAdmin)

	if out.Error != nil {
		return newAdmin, out.Error
	} else if out.RowsAffected == 0 {
		return newAdmin, gorm.ErrRecordNotFound
	}

	return newAdmin, nil

}

func (a AdminPostgresRepo) UpdateAdminRole(adminId string, newRole enums.AdminRole) error {

	out := a.db.
		Table("admins").
		Where("id=? AND deleted_at IS NULL AND (admin_role = 'member' OR admin_role = 'volunteer')", adminId).
		Update("admin_role", newRole)

	if out.Error != nil {
		return out.Error
	} else if out.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
