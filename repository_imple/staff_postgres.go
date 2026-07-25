package repositoryimple

import (
	"espectro/entity"
	"time"

	"gorm.io/gorm"
)

type StaffPostgresRepo struct {
	db *gorm.DB
}

func NewStaffPostgresRepo(db *gorm.DB) StaffPostgresRepo {
	return StaffPostgresRepo{db: db}
}

func (s StaffPostgresRepo) CreateStaff(staff entity.StaffEntity) (entity.StaffEntity, error) {
	err := s.db.
		Table("staffs").
		Create(&staff).Error
	return staff, err
}

func (s StaffPostgresRepo) UpdateStaff(staff entity.StaffUpdateEntity) (entity.StaffEntity, error) {
	var newStaff entity.StaffEntity

	out := s.db.Raw(
		`
		UPDATE staffs SET
		fullname=COALESCE(?,fullname),
		role=COALESCE(?,role),
		phone=COALESCE(?,phone),
		email=COALESCE(?,email)

		WHERE id=? AND deleted_at IS NULL
		RETURNING id,fullname,role,phone,email,created_at
		`,
		staff.Name, staff.Role, staff.PhoneNumber, staff.Email, staff.StaffId,
	).Scan(&newStaff)

	if out.Error != nil {
		return entity.StaffEntity{}, out.Error
	} else if out.RowsAffected == 0 {
		return entity.StaffEntity{}, gorm.ErrRecordNotFound
	}

	return newStaff, nil
}

func (s StaffPostgresRepo) DeleteStaff(staffId string) error {

	out := s.db.
		Table("staffs").
		Where("id=? AND deleted_at IS NULL", staffId).
		Update("deleted_at", time.Now().UTC())

	if out.Error != nil {
		return out.Error
	} else if out.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
