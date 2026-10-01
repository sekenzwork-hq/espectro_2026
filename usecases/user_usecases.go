package usecases

import (
	"errors"
	customerrors "espectro/custom_errors"
	"espectro/entity"
	"espectro/pkg"
	"espectro/repository"
	"strings"

	"github.com/bytedance/gopkg/util/logger"
	"gorm.io/gorm"
)

type UserUsecases struct {
	repo repository.UserRepository
}

func NewUserUsecases(r repository.UserRepository) UserUsecases {
	return UserUsecases{
		repo: r,
	}
}

func (u UserUsecases) RegisterUser(user entity.UserJsonCreateEntity) (entity.UserEntity, error) {

	empty := entity.UserEntity{}

	if len(user.Username) < 3 {
		return empty, &customerrors.ValidationError{DisplayError: "Username length should be greater than or equal to 3"}
	} else if len(user.Username) > 100 {
		return empty, &customerrors.ValidationError{DisplayError: "Username lengths should be less than or equal to 100"}
	}

	passwordErr := pkg.ValidatePassword(user.Password)

	if passwordErr != nil {
		return empty, passwordErr
	}

	fullnameErr := pkg.ValidateFullname(user.Fullname)

	if fullnameErr != nil {
		return empty, &customerrors.ValidationError{DisplayError: fullnameErr.Error()}
	}

	if !pkg.ValidateEmail(user.Email) {
		return empty, &customerrors.ValidationError{DisplayError: "Invalid email"}
	}

	if !pkg.ValidateCountryCode(user.CountryCode) {
		return empty, &customerrors.ValidationError{DisplayError: "Invalid country code"}
	}

	if !pkg.ValidateCountryOrState(user.Country) {
		return empty, &customerrors.ValidationError{DisplayError: "Invalid country name"}
	}

	if !pkg.ValidateCountryOrState(user.State) {
		return empty, &customerrors.ValidationError{DisplayError: "Invalid state name"}
	}

	if !pkg.ValidateCity(user.City) {
		return empty, &customerrors.ValidationError{DisplayError: "Invalid city"}
	}

	if !pkg.ValidatePhoneNumber(user.PhoneNumber) {
		return empty, &customerrors.ValidationError{DisplayError: "Invalid phone number"}
	}

	if !pkg.ValidateUserType(user.Usertype) {
		return empty, &customerrors.ValidationError{DisplayError: "Invalid user type"}
	}

	trimmedUsername := strings.TrimSpace(user.Username)
	trimmedFullname := strings.TrimSpace(user.Fullname)
	trimmedPassword := strings.TrimSpace(user.Password)
	trimmedCountry := strings.TrimSpace(user.Country)
	trimmedCountryCode := strings.TrimSpace(user.CountryCode)
	trimmedState := strings.TrimSpace(user.State)
	trimmedUserType := strings.TrimSpace(user.Usertype)
	trimmedCity := strings.TrimSpace(user.City)
	trimmedEmail := strings.TrimSpace(user.Email)
	trimmedPhoneNumber := strings.TrimSpace(user.PhoneNumber)

	exists, usernameCheckingErr := u.repo.UsernameExists(trimmedUsername)

	if usernameCheckingErr != nil {
		return empty, &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
	} else if exists {
		return empty, &customerrors.ValidationError{DisplayError: "Username already exists"}
	}

	hashedPass, hashErr := pkg.EncryptPassword(trimmedPassword)

	if hashErr != nil {
		return empty, &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
	}

	createdUser, dbErr := u.repo.RegisterUser(entity.UserDBCreateEntity{
		Fullname:    trimmedFullname,
		Username:    trimmedUsername,
		Password:    hashedPass,
		Email:       trimmedEmail,
		Country:     trimmedCountry,
		State:       trimmedState,
		City:        trimmedCity,
		PhoneNumber: trimmedPhoneNumber,
		CountryCode: trimmedCountryCode,
		Usertype:    trimmedUserType,
	})

	if dbErr != nil {
		logger.Error("DB error while inserting user data : ", dbErr)
		return empty, &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
	}

	return createdUser, nil
}

func (u UserUsecases) Login(userEnteredCred entity.UserCredentialsJsonEntity) (string, error) {

	credErr := &customerrors.CredentialsError{DisplayError: "Invalid credentials"}
	if len(userEnteredCred.Username) < 5 {
		return "", credErr
	} else if len(userEnteredCred.Username) > 100 {
		return "", credErr
	}

	if len(userEnteredCred.Password) < 6 {
		return "", credErr
	} else if len(userEnteredCred.Password) > 18 {
		return "", credErr
	}

	credentials, err := u.repo.RetrieveUserCredByUsername(userEnteredCred.Username)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", credErr
		} else {
			return "", &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
		}
	}

	if !pkg.CompareHashedPass(userEnteredCred.Password, credentials.Password) {
		return "", credErr
	}

	return credentials.Id, nil
}

func (u UserUsecases) CheckUserExists(userId string) error {

	if !pkg.ValidateUUID(userId) {
		return &customerrors.ValidationError{DisplayError: "Invalid user id"}
	}

	exists, err := u.repo.UserExists(userId)

	if err != nil {
		return &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
	}

	if !exists {
		return &customerrors.NotFoundError{DisplayError: "User does not exist"}
	}

	return nil
}
