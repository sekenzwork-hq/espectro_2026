package usecases

import (
	"errors"
	"espectro/entity"
	"espectro/pkg"
	"espectro/repository"

	"github.com/bytedance/gopkg/util/logger"
)

type UserUsecases struct {
	repo repository.UserPostgresRepo
}

func NewUserUsecases(r repository.UserPostgresRepo) UserUsecases {
	return UserUsecases{
		repo: r,
	}
}

func (u *UserUsecases) RegisterUser(user entity.UserEntity) (string, error) {

	fullnameErr := pkg.ValidateUserFullname(user.Fullname)

	if fullnameErr != nil {
		return "", fullnameErr
	}

	isEmailCorrect := pkg.ValidateUserEmail(user.Email)

	if !isEmailCorrect {
		return "", errors.New("Invalid Email address")
	}

	isCountryCodeCorrect := pkg.ValidateCountryCode(user.CountryCode)

	if !isCountryCodeCorrect {
		return "", errors.New("Invalid Country code")
	}

	isCountryCorrect := pkg.ValidateCountryOrState(user.Country)

	if !isCountryCorrect {
		return "", errors.New("Invalid Country name")
	}

	isStateCorrect := pkg.ValidateCountryOrState(user.State)

	if !isStateCorrect {
		return "", errors.New("Invalid State name")
	}

	isPhoneNumberCorrect := pkg.ValidatePhoneNumber(user.PhoneNumber)

	if !isPhoneNumberCorrect {
		return "", errors.New("Invalid Phone number")
	}

	isCorrectUserType := pkg.ValidateUserType(user.Usertype)

	if !isCorrectUserType {
		return "", errors.New("Invalid Usertype")
	}

	id, dbErr := u.RegisterUser(user)

	if dbErr != nil {
		logger.Error("DB error while inserting user data : ", dbErr)
		return "", errors.New("Something went wrong while saving user data")
	}

	return id, nil
}
