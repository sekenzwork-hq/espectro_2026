package usecases

import (
	customerrors "espectro/custom_errors"
	"espectro/entity"
	"espectro/pkg"
	"espectro/repository"
	"fmt"
)

type StaffUsecases struct {
	staffRepo repository.StaffRepo
}

func NewStaffUsecases(staffRepo repository.StaffRepo) StaffUsecases {
	return StaffUsecases{staffRepo: staffRepo}
}

func (e StaffUsecases) CreateStaff(staff entity.StaffFromJsonEntity) (entity.StaffEntity, error) {

	fmt.Println(staff.Name)
	emptyEntity := entity.StaffEntity{}
	nameErr := pkg.ValidateFullname(staff.Name)

	if nameErr != nil {
		return emptyEntity, &customerrors.ValidationError{OrgError: nameErr.Error()}
	} else if !staff.Role.IsValid() {
		return emptyEntity, &customerrors.ValidationError{OrgError: "Invalid role"}
	} else if !pkg.ValidateEmail(staff.Email) {
		return emptyEntity, &customerrors.ValidationError{OrgError: "Invalid email"}
	} else if !pkg.ValidatePhoneNumber(staff.PhoneNumber) {
		return emptyEntity, &customerrors.ValidationError{OrgError: "Invalid phone number"}
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
