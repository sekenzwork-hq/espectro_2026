package repository

import "espectro/entity"

type StaffRepo interface {
	CreateStaff(staff entity.StaffEntity) (entity.StaffEntity, error)
}
