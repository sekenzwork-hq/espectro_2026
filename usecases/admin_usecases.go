package usecases

import (
	"errors"
	customerrors "espectro/custom_errors"
	"espectro/entity"
	"espectro/pkg"
	"espectro/repository"

	"gorm.io/gorm"
)

type AdminUsecases struct {
	repo repository.AdminRepo
}

func NewAdminUsecases(repo repository.AdminRepo) AdminUsecases {
	return AdminUsecases{
		repo: repo,
	}
}

func (a AdminUsecases) Login(email string, password string) (string, error) {

	credentialError := &customerrors.CredentialsError{OrgError: "Invalid Credentials"}
	serverError := &customerrors.ServerError{OrgError: "Something went wrong while operating"}

	isEmailCorrect := pkg.ValidateEmail(email)

	if !isEmailCorrect || len(password) > 50 {
		return "", credentialError
	}

	adminCred, adminCredErr := a.repo.RetrieveAdminCredByEmail(email)

	if errors.Is(adminCredErr, gorm.ErrRecordNotFound) {
		return "", credentialError
	} else if adminCredErr != nil {
		return "", serverError
	}

	hashedPass := adminCred.Password

	passCorrect := pkg.CompareHashedPass(password, hashedPass)

	if !passCorrect {
		return "", credentialError
	}

	jwtToken, jwtErr := pkg.GenerateJWTForAdmin(adminCred.Id.String())

	if jwtErr != nil {
		return "", serverError
	}

	return jwtToken, nil

}

func (a AdminUsecases) CreateNewAdmin(admin entity.AdminEntity, requestedAdminId string) (string, error) {

	if len(requestedAdminId) == 0 {
		return "", &customerrors.CredentialsError{OrgError: "Current admin is invalid"}
	}

	//First retrieving current admin role to detect whether the admin has the permission to do it.
	currentAdminRole, currentAdminRoleErr := a.repo.RetrieveAdminRoleByID(requestedAdminId)

	if currentAdminRoleErr != nil {
		return "", &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}

	if len(currentAdminRole) == 0 {
		return "", &customerrors.CredentialsError{OrgError: "Current admin is invalid"}
	}

	if currentAdminRole != "leader" {
		return "", &customerrors.CredentialsError{OrgError: "Current admin doesn't have permission"}
	}

	fullnameErr := pkg.ValidateFullname(admin.Fullname)

	if fullnameErr != nil {
		return "", &customerrors.ValidationError{OrgError: fullnameErr.Error()}
	}

	isEmailCorrect := pkg.ValidateEmail(admin.Email)

	if !isEmailCorrect {
		return "", &customerrors.ValidationError{OrgError: "Invalid email address"}
	}

	passwordErr := pkg.ValidatePassword(admin.Password)

	if passwordErr != nil {
		return "", &customerrors.ValidationError{OrgError: passwordErr.Error()}
	}

	isRoleCorrect := pkg.ValidateAdminRole(admin.Role)

	if !isRoleCorrect {
		return "", &customerrors.ValidationError{OrgError: "Invalid admin role"}
	}

	hashedPass, hashingErr := pkg.EncryptPassword(admin.Password)

	if hashingErr != nil {
		return "", &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}

	adminWithPasswordHashed := entity.AdminEntity{
		Fullname: admin.Fullname,
		Email:    admin.Email,
		Role:     admin.Role,
		Password: hashedPass,
	}

	newAdminId, insertErr := a.repo.CreateNewAdmin(adminWithPasswordHashed)

	if insertErr != nil {
		return "", &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}

	return newAdminId.String(), nil

}

func (a AdminUsecases) DeleteMemberOrVolunteer(adminId string, requestedAdminId string) (bool, error) {

	if len(requestedAdminId) == 0 {
		return false, &customerrors.CredentialsError{OrgError: "Current admin is invalid"}
	} else if len(adminId) == 0 {
		return false, &customerrors.ValidationError{OrgError: "Provide valid admin id to delete"}
	}

	//First checking the current admin role to check the permission.
	currentAdminRole, currentAdminRoleErr := a.repo.RetrieveAdminRoleByID(requestedAdminId)

	if currentAdminRoleErr != nil {
		return false, &customerrors.CredentialsError{OrgError: "Something went wrong while operating"}
	}

	if currentAdminRole != "leader" {
		return false, &customerrors.CredentialsError{OrgError: "Current admin doesn't have permission"}
	}

	deletionErr := a.repo.DeleteMemberOrVolunteer(adminId)

	if deletionErr != nil {
		return false, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}

	return true, nil
}
