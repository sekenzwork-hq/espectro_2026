package repositoryimple

import (
	customerrors "espectro/custom_errors"
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

func (a AdminPostgresRepo) CreateNewAdmin(admin entity.AdminEntity) (uuid.UUID, error) {
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

func (a AdminPostgresRepo) UpdateCurrentAdmin(
	adminId string,
	newEmail string,
	newFullname string,
	updateMode enums.AdminUpdateMode) error {

	data := map[string]any{}

	switch updateMode {

	case enums.EmailOnly:
		data["email"] = newEmail

	case enums.FullnameOnly:
		data["fullname"] = newFullname

	case enums.EmailAndFullname:
		data["email"] = newEmail
		data["fullname"] = newFullname

	default:
		return &customerrors.InvalidAdminUpdateModeError{OrgError: "Provide correct update mode"}

	}

	out := a.db.
		Table("admins").
		Where("id=? AND deleted_at IS NULL", adminId).
		Updates(data)

	if out.Error != nil {
		return out.Error
	} else if out.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil

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
