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

func (a AdminUsecases) CreateNewAdmin(admin entity.AdminEntity, requestedAdminId string) (string, error) {

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

func (a AdminUsecases) UpdateCurrentAdmin(adminId string, newEmail string, newFullname string) error {

	if len(adminId) == 0 {
		return &customerrors.AuthenticationError{OrgError: "Current admin is invalid"}
	}

	var updationMode enums.AdminUpdateMode

	if len(newEmail) != 0 && len(newFullname) != 0 {
		updationMode = enums.EmailAndFullname
	} else if len(newEmail) != 0 {
		updationMode = enums.EmailOnly
	} else {
		updationMode = enums.FullnameOnly
	}

	var outErr error

	updateFunc := func(newEmail string, newFullname string, updateMode enums.AdminUpdateMode) error {

		err := a.repo.UpdateCurrentAdmin(adminId, newEmail, newFullname, updateMode)

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &customerrors.NotFoundError{OrgError: "Current Admin is invalid"}
		} else if err != nil {
			return &customerrors.ServerError{OrgError: "Something went wrong while operating"}
		}

		return nil

	}
	switch updationMode {

	//First checking if the request is to update fullname and email
	case enums.EmailAndFullname:
		isEmailCorrect := pkg.ValidateEmail(newEmail)
		fullnameErr := pkg.ValidateFullname(newFullname)

		if !isEmailCorrect {
			return &customerrors.ValidationError{OrgError: "Invalid email address"}
		} else if fullnameErr != nil {
			return &customerrors.ValidationError{OrgError: fullnameErr.Error()}
		}
		err := updateFunc(newEmail, newFullname, enums.EmailAndFullname)
		outErr = err

	//Checking if the request is to update fullname only
	case enums.FullnameOnly:

		fullnameErr := pkg.ValidateFullname(newFullname)

		if fullnameErr != nil {
			return &customerrors.ValidationError{OrgError: fullnameErr.Error()}
		}
		err := updateFunc("", newFullname, enums.FullnameOnly)
		outErr = err

	//Checking if the request is to update email only
	case enums.EmailOnly:
		isEmailCorrect := pkg.ValidateEmail(newEmail)

		if !isEmailCorrect {
			return &customerrors.ValidationError{OrgError: "Invalid email address"}
		}
		err := updateFunc(newEmail, "", enums.EmailOnly)
		outErr = err

	//If the above conditions are not satisfied, then [updationMode] may be not valid, so just returning validation error
	default:
		outErr = &customerrors.ValidationError{OrgError: "Provide email, fullname or both to update"}
	}

	return outErr

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
