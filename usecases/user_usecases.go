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
	repo repository.UserRepository
}

func NewUserUsecases(r repository.UserRepository) UserUsecases {
	return UserUsecases{
		repo: r,
	}
}

func (u *UserUsecases) RegisterUser(user entity.UserEntity) (uuid.UUID, error) {

	fullnameErr := pkg.ValidateFullname(user.Fullname)

	if fullnameErr != nil {
		return uuid.UUID{}, &customerrors.ValidationError{OrgError: fullnameErr.Error()}
	}

	isEmailCorrect := pkg.ValidateEmail(user.Email)

	if !isEmailCorrect {
		return uuid.UUID{}, &customerrors.ValidationError{OrgError: "Invalid email"}
	}

	isCountryCodeCorrect := pkg.ValidateCountryCode(user.CountryCode)

	if !isCountryCodeCorrect {
		return uuid.UUID{}, &customerrors.ValidationError{OrgError: "Invalid country code"}
	}

	isCountryCorrect := pkg.ValidateCountryOrState(user.Country)

	if !isCountryCorrect {
		return uuid.UUID{}, &customerrors.ValidationError{OrgError: "Invalid country name"}
	}

	isStateCorrect := pkg.ValidateCountryOrState(user.State)

	if !isStateCorrect {
		return uuid.UUID{}, &customerrors.ValidationError{OrgError: "Invalid state name"}
	}

	isCityCorrect := pkg.ValidateCity(user.City)

	if !isCityCorrect {
		return uuid.UUID{}, &customerrors.ValidationError{OrgError: "Invalid city"}
	}

	isPhoneNumberCorrect := pkg.ValidatePhoneNumber(user.PhoneNumber)

	if !isPhoneNumberCorrect {
		return uuid.UUID{}, &customerrors.ValidationError{OrgError: "Invalid phone number"}
	}

	isCorrectUserType := pkg.ValidateUserType(user.Usertype)

	if !isCorrectUserType {
		return uuid.UUID{}, &customerrors.ValidationError{OrgError: "Invalid user type"}
	}

	id, dbErr := u.repo.RegisterUser(user)

	if dbErr != nil {
		logger.Error("DB error while inserting user data : ", dbErr)
		return uuid.UUID{}, &customerrors.ServerError{OrgError: "Something went wrong while saving user data"}
	}

	return id, nil
}
