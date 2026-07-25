package repository

import "espectro/entity"

type StaffRepo interface {
	CreateStaff(staff entity.StaffEntity) (entity.StaffEntity, error)
	UpdateStaff(staff entity.StaffUpdateEntity) (entity.StaffEntity, error)
	DeleteStaff(staffId string) error
}
