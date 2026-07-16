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

func (a AdminUsecases) CreateNewAdmin(admin entity.AdminCreateEntity, requestedAdminId string) (string, error) {

	var adminRole enums.AdminRole

	fullnameErr := pkg.ValidateFullname(admin.Fullname)
	adminRole, isAdminRoleCorrect := adminRole.ParseRole(admin.Role)
	isEmailCorrect := pkg.ValidateEmail(admin.Email)
	passwordErr := pkg.ValidatePassword(admin.Password)

	if fullnameErr != nil {
		return "", &customerrors.ValidationError{OrgError: fullnameErr.Error()}
	} else if !isEmailCorrect {
		return "", &customerrors.ValidationError{OrgError: "Invalid email address"}
	} else if passwordErr != nil {
		return "", &customerrors.ValidationError{OrgError: passwordErr.Error()}
	} else if !isAdminRoleCorrect {
		return "", &customerrors.ValidationError{OrgError: "Invalid admin role"}
	}

	hashedPass, hashingErr := pkg.EncryptPassword(admin.Password)
	if hashingErr != nil {
		return "", &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}

	adminWithPasswordHashed := entity.AdminCreateEntity{
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

func (a AdminUsecases) DeleteMemberOrVolunteer(adminId string, requestedAdminId string) error {

	if len(adminId) == 0 {
		return &customerrors.ValidationError{OrgError: "Provide valid admin id to delete"}
	}

	deletionErr := a.repo.DeleteMemberOrVolunteer(adminId)
	if errors.Is(deletionErr, gorm.ErrRecordNotFound) {
		return &customerrors.NotFoundOrLeaderError{OrgError: "Deletion operation doesn't work whether admin doesn't exist or admin is a leader"}
	} else if deletionErr != nil {
		return &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}

	return nil
}

func (a AdminUsecases) CheckAdminExists(adminId string) error {

	if len(adminId) == 0 {
		return &customerrors.AuthenticationError{OrgError: "Current admin is invalid"}
	}

	exists, err := a.repo.CheckAdminExists(adminId)
	if err != nil {
		return &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	} else if !exists {
		return &customerrors.AuthenticationError{OrgError: "Current admin is invalid"}
	}

	return nil
}

func (a AdminUsecases) UpdateCurrentAdmin(adminId string, newAdmin entity.AdminUpdateEntity) (entity.AdminEntity, error) {

	emptyAdmin := entity.AdminEntity{}
	if !pkg.ValidateUUID(adminId) {
		return emptyAdmin, &customerrors.AuthenticationError{OrgError: "Admin does not exist"}
	}

	updatedAdmin, err := a.repo.UpdateCurrentAdmin(adminId, newAdmin)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return emptyAdmin, &customerrors.NotFoundError{OrgError: "Current Admin is invalid"}
	} else if err != nil {
		return emptyAdmin, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}
	return updatedAdmin, nil

}

func (a AdminUsecases) RetrieveAdminRoleByID(adminId string) (enums.AdminRole, error) {

	if adminId == "" {
		return "", &customerrors.AuthenticationError{OrgError: "Current admin is invalid"}
	}

	role, err := a.repo.RetrieveAdminRoleByID(adminId)

	if errors.Is(err, gorm.ErrRecordNotFound) || role == "" {
		return "", &customerrors.AuthenticationError{OrgError: "Current admin is invalid"}
	} else if err != nil {
		return "", &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	} else {
		return role, nil
	}

}
func (a AdminUsecases) UpdateAdminRole(adminId string, newRole string) error {

	if len(adminId) == 0 || len(adminId) > 36 {
		return &customerrors.ValidationError{OrgError: "Provide valid admin id"}
	}

	var adminRole enums.AdminRole
	adminRole, isRoleCorrect := adminRole.ParseRole(newRole)

	if !isRoleCorrect {
		return &customerrors.ValidationError{OrgError: "Invalid admin role"}
	}

	err := a.repo.UpdateAdminRole(adminId, adminRole)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &customerrors.NotFoundOrLeaderError{OrgError: "Updation operation doesn't work whether the admin doesn't exist or admin is a leader"}
	} else if err != nil {
		return &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}

	return nil

}
