package usecases

import (
	customerrors "espectro/custom_errors"
	"espectro/entity"
	"espectro/pkg"
	"espectro/repository"
	"mime/multipart"

	"github.com/google/uuid"
)

type InvestorUsecases struct {
	investorRepo repository.InvestorRepo
	mediaRepo    repository.MediaServiceRepo
}

func NewInvestorUsecases(investorRepo repository.InvestorRepo, mediaRepo repository.MediaServiceRepo) InvestorUsecases {
	return InvestorUsecases{investorRepo: investorRepo, mediaRepo: mediaRepo}
}

func (i InvestorUsecases) AddInvestor(name string, websiteUrl *string, phoneNumber string, email string, logo *multipart.FileHeader) (entity.InvestorEntity, error) {

	emptyInvestor := entity.InvestorEntity{}

	if err := pkg.ValidateName(name); err != nil {
		return emptyInvestor, &customerrors.ValidationError{OrgError: err.Error()}
	}

	if !pkg.ValidatePhoneNumber(phoneNumber) {
		return emptyInvestor, &customerrors.ValidationError{OrgError: "Invalid phone number"}
	}

	if !pkg.ValidateEmail(email) {
		return emptyInvestor, &customerrors.ValidationError{OrgError: "Invalid email"}
	}

	var logoUrl *string
	investorId := uuid.New().String()

	if logo != nil {
		url, uploadErr := i.mediaRepo.UploadFile(logo, "investor/"+investorId)
		if uploadErr != nil {
			return emptyInvestor, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
		}
		logoUrl = &url
	}

	investorToInsert := entity.InvestorEntity{
		WebsiteUrl:  websiteUrl,
		PhoneNumber: phoneNumber,
		Email:       email,
		LogoUrl:     logoUrl,
		Id:          investorId,
	}

	newInvestor, insertionErr := i.investorRepo.AddInvestor(investorToInsert)

	if insertionErr != nil {
		i.mediaRepo.DeleteFolderWithFiles("investor/", investorId)
		return emptyInvestor, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}

	return newInvestor, nil

}
