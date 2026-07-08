package usecases

import (
	customerrors "espectro/custom_errors"
	"espectro/models"
	"espectro/pkg"
	"espectro/repository"
)

type AdminUsecases struct {
	repo repository.AdminRepo
}

func NewAdminUsecases(repo repository.AdminRepo) AdminUsecases {
	return AdminUsecases{
		repo: repo,
	}
}

func (a AdminUsecases) RetrieveAdminCredByEmail(email string, password string) (models.TokensModel, error) {

	isEmailCorrect := pkg.ValidateUserEmail(email)

	if !isEmailCorrect {
		return models.TokensModel{}, &customerrors.CredentialsError{OrgError: "Invalid Credentials"}
	}

	adminCred, adminCredErr := a.repo.RetrieveAdminCredByEmail(email)

	if adminCredErr != nil {
		return models.TokensModel{}, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}

	hashedPass := adminCred.Password

	passCorrect := pkg.CompareHashedPass(password, hashedPass)

	if !passCorrect {
		return models.TokensModel{}, &customerrors.CredentialsError{OrgError: "Invalid Credentials"}
	}

}
