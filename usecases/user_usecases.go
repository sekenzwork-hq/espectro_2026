package usecases

import (
	customerrors "espectro/custom_errors"
	"espectro/entity"
	"espectro/pkg"
	"espectro/repository"

	"github.com/bytedance/gopkg/util/logger"
)

type UserUsecases struct {
	repo repository.UserRepository
}

func NewUserUsecases(r repository.UserRepository) UserUsecases {
	return UserUsecases{
		repo: r,
	}
}

func (u UserUsecases) RegisterUser(user entity.UserEntity) (entity.UserEntity, error) {

	empty := entity.UserEntity{}
	fullnameErr := pkg.ValidateFullname(user.Fullname)

	if fullnameErr != nil {
		return empty, &customerrors.ValidationError{DisplayError: fullnameErr.Error()}
	}

	isEmailCorrect := pkg.ValidateEmail(user.Email)

	if !isEmailCorrect {
		return empty, &customerrors.ValidationError{DisplayError: "Invalid email"}
	}

	isCountryCodeCorrect := pkg.ValidateCountryCode(user.CountryCode)

	if !isCountryCodeCorrect {
		return empty, &customerrors.ValidationError{DisplayError: "Invalid country code"}
	}

	isCountryCorrect := pkg.ValidateCountryOrState(user.Country)

	if !isCountryCorrect {
		return empty, &customerrors.ValidationError{DisplayError: "Invalid country name"}
	}

	isStateCorrect := pkg.ValidateCountryOrState(user.State)

	if !isStateCorrect {
		return empty, &customerrors.ValidationError{DisplayError: "Invalid state name"}
	}

	isCityCorrect := pkg.ValidateCity(user.City)

	if !isCityCorrect {
		return empty, &customerrors.ValidationError{DisplayError: "Invalid city"}
	}

	isPhoneNumberCorrect := pkg.ValidatePhoneNumber(user.PhoneNumber)

	if !isPhoneNumberCorrect {
		return empty, &customerrors.ValidationError{DisplayError: "Invalid phone number"}
	}

	isCorrectUserType := pkg.ValidateUserType(user.Usertype)

	if !isCorrectUserType {
		return empty, &customerrors.ValidationError{DisplayError: "Invalid user type"}
	}

	createdUser, dbErr := u.repo.RegisterUser(user)

	if dbErr != nil {
		logger.Error("DB error while inserting user data : ", dbErr)
		return empty, &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
	}

	return createdUser, nil
}

func (u UserUsecases) CheckUserExists(userId string) (bool, error) {

	if !pkg.ValidateUUID(userId) {
		return false, &customerrors.ValidationError{DisplayError: "Invalid user id"}
	}

	exists, err := u.repo.CheckUserExists(userId)

	if err != nil {
		return false, &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
	}

	return exists, nil
}
