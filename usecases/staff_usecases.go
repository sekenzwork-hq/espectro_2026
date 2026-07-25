package usecases

import (
	"errors"
	customerrors "espectro/custom_errors"
	"espectro/entity"
	"espectro/enums"
	"espectro/pkg"
	"espectro/repository"

	"gorm.io/gorm"
)

type StaffUsecases struct {
	staffRepo repository.StaffRepo
}

func NewStaffUsecases(staffRepo repository.StaffRepo) StaffUsecases {
	return StaffUsecases{staffRepo: staffRepo}
}

func (e StaffUsecases) CreateStaff(staff entity.StaffFromJsonEntity) (entity.StaffEntity, error) {

	emptyEntity := entity.StaffEntity{}

	validationErr := e.validateStaffData(nil, &staff.Name, &staff.Role, &staff.PhoneNumber, &staff.Email)
	if validationErr != nil {
		return emptyEntity, validationErr
	}
	newStaff, err := e.staffRepo.CreateStaff(entity.StaffEntity{
		Name:        staff.Name,
		Role:        staff.Role,
		PhoneNumber: staff.PhoneNumber,
		Email:       staff.Email,
	})

	if err != nil {
		return emptyEntity, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}

	return newStaff, nil

}

func (e StaffUsecases) UpdateStaff(staff entity.StaffUpdateEntity) (entity.StaffEntity, error) {

	empty := entity.StaffEntity{}
	validationErr := e.validateStaffData(&staff.StaffId, staff.Name, staff.Role, staff.PhoneNumber, staff.Email)
	if validationErr != nil {
		return empty, validationErr
	}

	newStaff, updationErr := e.staffRepo.UpdateStaff(staff)

	if errors.Is(updationErr, gorm.ErrRecordNotFound) {
		return empty, &customerrors.NotFoundError{OrgError: "Staff does not exist"}
	} else if updationErr != nil {
		return empty, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}

	return newStaff, nil
}

func (e StaffUsecases) DeleteStaff(staffId string) error {

	if !pkg.ValidateUUID(staffId) {
		return &customerrors.ValidationError{OrgError: "Invalid staff id"}
	}

	err := e.staffRepo.DeleteStaff(staffId)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &customerrors.NotFoundError{OrgError: "Staff does not exist"}
	} else if err != nil {
		return &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}

	return nil
}
func (e StaffUsecases) validateStaffData(staffId *string, fullname *string, role *enums.StaffRole, phoneNumber *string, email *string) error {

	if staffId != nil && !pkg.ValidateUUID(*staffId) {
		return &customerrors.ValidationError{OrgError: "Invalid staff id"}
	}
	if fullname != nil {
		nameErr := pkg.ValidateFullname(*fullname)

		if nameErr != nil {
			return &customerrors.ValidationError{OrgError: nameErr.Error()}
		}
	}

	if role != nil && !role.IsValid() {
		return &customerrors.ValidationError{OrgError: "Invalid role"}
	} else if email != nil && !pkg.ValidateEmail(*email) {
		return &customerrors.ValidationError{OrgError: "Invalid email"}
	} else if phoneNumber != nil && !pkg.ValidatePhoneNumber(*phoneNumber) {
		return &customerrors.ValidationError{OrgError: "Invalid phone number"}
	}

	return nil
}
