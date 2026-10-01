package usecases

import (
	"errors"
	customerrors "espectro/custom_errors"
	"espectro/entity"
	"espectro/enums"
	"espectro/pkg"
	"espectro/repository"
	"strings"

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

	credentialError := &customerrors.CredentialsError{DisplayError: "Invalid Credentials"}
	serverError := &customerrors.ServerError{DisplayError: "Something went wrong while operating"}

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

func (a AdminUsecases) CreateNewAdmin(admin entity.AdminCreateEntity, requestedAdminId string) (entity.AdminEntity, error) {

	empty := entity.AdminEntity{}

	var adminRole enums.AdminRole

	fullnameErr := pkg.ValidateFullname(admin.Fullname)
	adminRole, isAdminRoleCorrect := adminRole.ParseRole(admin.Role)
	isEmailCorrect := pkg.ValidateEmail(admin.Email)
	passwordErr := pkg.ValidatePassword(admin.Password)

	if fullnameErr != nil {
		return empty, &customerrors.ValidationError{DisplayError: fullnameErr.Error()}
	} else if !isEmailCorrect {
		return empty, &customerrors.ValidationError{DisplayError: "Invalid email"}
	} else if passwordErr != nil {
		return empty, &customerrors.ValidationError{DisplayError: passwordErr.Error()}
	} else if !isAdminRoleCorrect {
		return empty, &customerrors.ValidationError{DisplayError: "Invalid admin role"}
	}

	trimmedEmail := strings.TrimSpace(admin.Email)

	emailExists, emailCheckingErr := a.repo.AdminEmailExists(trimmedEmail)

	if emailCheckingErr != nil {
		return empty, &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
	} else if emailExists {
		return empty, &customerrors.ValidationError{DisplayError: "Email already exists"}
	}

	hashedPass, hashingErr := pkg.EncryptPassword(admin.Password)
	if hashingErr != nil {
		return empty, &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
	}

	adminWithPasswordHashed := entity.AdminCreateEntity{
		Fullname: admin.Fullname,
		Email:    admin.Email,
		Role:     admin.Role,
		Password: hashedPass,
	}

	createdAdmin, insertErr := a.repo.CreateNewAdmin(adminWithPasswordHashed)

	if insertErr != nil {
		return empty, &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
	}

	return createdAdmin, nil

}

func (a AdminUsecases) DeleteMemberOrVolunteer(adminId string, requestedAdminId string) error {

	if len(adminId) == 0 {
		return &customerrors.ValidationError{DisplayError: "Provide valid admin id to delete"}
	}

	deletionErr := a.repo.DeleteMemberOrVolunteer(adminId)
	if errors.Is(deletionErr, gorm.ErrRecordNotFound) {
		return &customerrors.NotFoundOrLeaderError{DisplayError: "Deletion operation doesn't work whether admin doesn't exist or admin is a leader"}
	} else if deletionErr != nil {
		return &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
	}

	return nil
}

func (a AdminUsecases) AdminExists(adminId string) error {

	if len(adminId) == 0 {
		return &customerrors.AuthenticationError{DisplayError: "Current admin is invalid"}
	}

	exists, err := a.repo.AdminExists(adminId)
	if err != nil {
		return &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
	} else if !exists {
		return &customerrors.AuthenticationError{DisplayError: "Current admin is invalid"}
	}

	return nil
}

func (a AdminUsecases) UpdateCurrentAdmin(adminId string, newAdmin entity.AdminUpdateEntity) (entity.AdminEntity, error) {

	emptyAdmin := entity.AdminEntity{}
	if !pkg.ValidateUUID(adminId) {
		return emptyAdmin, &customerrors.AuthenticationError{DisplayError: "Admin does not exist"}
	}

	updatedAdmin, err := a.repo.UpdateCurrentAdmin(adminId, newAdmin)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return emptyAdmin, &customerrors.NotFoundError{DisplayError: "Current Admin is invalid"}
	} else if err != nil {
		return emptyAdmin, &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
	}
	return updatedAdmin, nil

}

func (a AdminUsecases) RetrieveAdminRoleByID(adminId string) (enums.AdminRole, error) {

	if adminId == "" {
		return "", &customerrors.AuthenticationError{DisplayError: "Current admin is invalid"}
	}

	role, err := a.repo.RetrieveAdminRoleByID(adminId)

	if errors.Is(err, gorm.ErrRecordNotFound) || role == "" {
		return "", &customerrors.AuthenticationError{DisplayError: "Current admin is invalid"}
	} else if err != nil {
		return "", &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
	} else {
		return role, nil
	}

}
func (a AdminUsecases) UpdateAdminRole(adminId string, newRole string) error {

	if len(adminId) == 0 || len(adminId) > 36 {
		return &customerrors.ValidationError{DisplayError: "Provide valid admin id"}
	}

	var adminRole enums.AdminRole
	adminRole, isRoleCorrect := adminRole.ParseRole(newRole)

	if !isRoleCorrect {
		return &customerrors.ValidationError{DisplayError: "Invalid admin role"}
	}

	err := a.repo.UpdateAdminRole(adminId, adminRole)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &customerrors.NotFoundOrLeaderError{DisplayError: "Updation operation doesn't work whether the admin doesn't exist or admin is a leader"}
	} else if err != nil {
		return &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
	}

	return nil

}
