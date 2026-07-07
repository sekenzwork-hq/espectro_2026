package usecases

import (
	customerrors "espectro/custom_errors"
	"espectro/entity"
	"espectro/pkg"
	"espectro/repository"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/google/uuid"
)

type UserUsecases struct {
	repo repository.UserPostgresRepo
}

func NewUserUsecases(r repository.UserPostgresRepo) UserUsecases {
	return UserUsecases{
		repo: r,
	}
}

func (u *UserUsecases) RegisterUser(user entity.UserEntity) (uuid.UUID, error) {

	fullnameErr := pkg.ValidateUserFullname(user.Fullname)

	if fullnameErr != nil {
		return uuid.UUID{}, &customerrors.ValidationError{OrgError: fullnameErr.Error()}
	}

	isEmailCorrect := pkg.ValidateUserEmail(user.Email)

	if !isEmailCorrect {
		return uuid.UUID{}, &customerrors.ValidationError{OrgError: "Invalid Email address"}
	}

	isCountryCodeCorrect := pkg.ValidateCountryCode(user.CountryCode)

	if !isCountryCodeCorrect {
		return uuid.UUID{}, &customerrors.ValidationError{OrgError: "Invalid Country code"}
	}

	isCountryCorrect := pkg.ValidateCountryOrState(user.Country)

	if !isCountryCorrect {
		return uuid.UUID{}, &customerrors.ValidationError{OrgError: "Invalid Country name"}
	}

	isStateCorrect := pkg.ValidateCountryOrState(user.State)

	if !isStateCorrect {
		return uuid.UUID{}, &customerrors.ValidationError{OrgError: "Invalid State name"}
	}

	isPhoneNumberCorrect := pkg.ValidatePhoneNumber(user.PhoneNumber)

	if !isPhoneNumberCorrect {
		return uuid.UUID{}, &customerrors.ValidationError{OrgError: "Invalid Phone number"}
	}

	isCorrectUserType := pkg.ValidateUserType(user.Usertype)

	if !isCorrectUserType {
		return uuid.UUID{}, &customerrors.ValidationError{OrgError: "Invalid Usertype"}
	}

	id, dbErr := u.repo.RegisterUser(user)

	if dbErr != nil {
		logger.Error("DB error while inserting user data : ", dbErr)
		return uuid.UUID{}, &customerrors.ServerError{OrgError: "Something went wrong while saving user data"}
	}

	return id, nil
}
