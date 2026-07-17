package usecases

import (
	"errors"
	customerrors "espectro/custom_errors"
	"espectro/entity"
	"espectro/pkg"
	"espectro/repository"
	"fmt"
	"mime/multipart"

	"github.com/google/uuid"
	"gorm.io/gorm"
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

	if err := i.validateInvestorDetails(nil, &name, &phoneNumber, &email, websiteUrl); err != nil {
		return emptyInvestor, err
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
		Name:        name,
	}

	newInvestor, insertionErr := i.investorRepo.AddInvestor(investorToInsert)
	if insertionErr != nil {
		i.mediaRepo.DeleteFolderWithFiles("investor/", investorId)
		return emptyInvestor, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}

	return newInvestor, nil

}

func (i InvestorUsecases) UpdateInvestor(investorId string, name *string, phoneNumber *string, email *string, websiteUrl *string, logo *multipart.FileHeader) (entity.InvestorEntity, error) {

	emptyInvestor := entity.InvestorEntity{}

	if err := i.validateInvestorDetails(&investorId, name, phoneNumber, email, websiteUrl); err != nil {
		fmt.Println("Err : ", err)
		return emptyInvestor, err
	}

	fmt.Println("Validation was correct")

	var logoUrl *string

	if logo != nil {
		url, err := i.mediaRepo.UploadFile(logo, "investor/"+investorId)
		if err != nil {
			return emptyInvestor, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
		}
		logoUrl = &url
	}

	investorToUpdate := entity.InvestorUpdateEntity{
		Name:        name,
		WebsiteUrl:  websiteUrl,
		PhoneNumber: phoneNumber,
		Email:       email,
		LogoUrl:     logoUrl,
	}

	newInvestor, updationErr := i.investorRepo.UpdateInvestor(investorId, investorToUpdate)

	if errors.Is(updationErr, gorm.ErrRecordNotFound) {
		i.mediaRepo.DeleteFolderWithFiles("investor/", investorId)
		return newInvestor, &customerrors.NotFoundError{OrgError: "Investor does not exist"}
	} else if updationErr != nil {
		i.mediaRepo.DeleteFolderWithFiles("investor/", investorId)
		return newInvestor, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}

	return newInvestor, nil

}

func (i InvestorUsecases) validateInvestorDetails(id *string, name *string, phoneNumber *string, email *string, websiteUrl *string) error {
	if id != nil && !pkg.ValidateUUID(*id) {
		return &customerrors.ValidationError{OrgError: "Invalid investor id"}
	}
	if name != nil {
		if err := pkg.ValidateName(*name); err != nil {
			return &customerrors.ValidationError{OrgError: err.Error()}
		}
	}
	if phoneNumber != nil && !pkg.ValidatePhoneNumber(*phoneNumber) {
		return &customerrors.ValidationError{OrgError: "Invalid phone number"}
	}
	if email != nil && !pkg.ValidateEmail(*email) {
		return &customerrors.ValidationError{OrgError: "Invalid email"}
	}
	if websiteUrl != nil {
		if err := pkg.ValidateUrl(*websiteUrl, "website url"); err != nil {
			return err
		}
	}
	return nil
}
