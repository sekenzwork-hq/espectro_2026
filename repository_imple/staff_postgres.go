package repositoryimple

import (
	"espectro/entity"

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
