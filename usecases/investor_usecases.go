package usecases

import (
	"errors"
	customerrors "espectro/custom_errors"
	"espectro/entity"
	"espectro/pkg"
	"espectro/repository"
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

func (i InvestorUsecases) CreateInvestor(name string, websiteUrl *string, phoneNumber string, email string, logo *multipart.FileHeader) (entity.InvestorEntity, error) {

	emptyInvestor := entity.InvestorEntity{}

	if err := i.validateInvestorDetails(nil, &name, &phoneNumber, &email, websiteUrl, logo); err != nil {
		return emptyInvestor, err
	}

	var logoUrl *string
	investorId := uuid.NewString()

	if logo != nil {
		url, _, uploadErr := i.mediaRepo.UploadFile(logo, "investor/"+investorId, true)
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

	newInvestor, insertionErr := i.investorRepo.CreateInvestor(investorToInsert)
	if insertionErr != nil {
		i.mediaRepo.DeleteFile("investor/", investorId)
		return emptyInvestor, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}

	return newInvestor, nil

}

func (i InvestorUsecases) UpdateInvestor(investorId string, name *string, phoneNumber *string, email *string, websiteUrl *string, logo *multipart.FileHeader) (entity.InvestorEntity, error) {

	emptyInvestor := entity.InvestorEntity{}

	if err := i.validateInvestorDetails(&investorId, name, phoneNumber, email, websiteUrl, logo); err != nil {
		return emptyInvestor, err
	}

	var logoUrl *string
	var newPublicId string
	var oldPublicId string
	folderId := "investor/" + investorId
	if logo != nil {
		oldId, err := i.mediaRepo.RetrieveAssetPublicIds(folderId)
		if err != nil {
			return emptyInvestor, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
		}
		if len(oldId) != 0 {
			oldPublicId = oldId[0]
		}
		url, id, err := i.mediaRepo.UploadFile(logo, folderId, false)
		if err != nil {
			return emptyInvestor, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
		}
		logoUrl = &url
		newPublicId = id
	}

	investorToUpdate := entity.InvestorUpdateEntity{
		Name:        name,
		WebsiteUrl:  websiteUrl,
		PhoneNumber: phoneNumber,
		Email:       email,
		LogoUrl:     logoUrl,
	}

	newInvestor, updationErr := i.investorRepo.UpdateInvestor(investorId, investorToUpdate)

	go func() {
		if updationErr != nil {
			i.mediaRepo.DeleteAssetsWithPublicIds([]string{newPublicId})
		} else {
			i.mediaRepo.DeleteAssetsWithPublicIds([]string{oldPublicId})
		}
	}()

	if errors.Is(updationErr, gorm.ErrRecordNotFound) {
		return newInvestor, &customerrors.NotFoundError{OrgError: "Investor does not exist"}
	} else if updationErr != nil {
		return newInvestor, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}

	return newInvestor, nil

}

func (i InvestorUsecases) DeleteInvestor(investorId string) error {

	if !pkg.ValidateUUID(investorId) {
		return &customerrors.ValidationError{OrgError: "Invalid investor id"}
	}

	err := i.investorRepo.DeleteInvestor(investorId)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &customerrors.NotFoundError{OrgError: "Investor does not exist"}
	} else if err != nil {
		return &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}

	return nil

}

func (i InvestorUsecases) RetrieveInvestors(limit int, page int) ([]entity.InvestorEntity, error) {

	offset := pkg.GetOffset(limit, page)
	investors, err := i.investorRepo.RetrieveInvestors(limit, offset)
	if err != nil {
		return investors, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}

	return investors, nil
}

func (i InvestorUsecases) validateInvestorDetails(id *string, name *string, phoneNumber *string, email *string, websiteUrl *string, logo *multipart.FileHeader) error {

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
	if logo != nil {
		correct, err := pkg.ValidateImage(logo)
		if err != nil {
			return &customerrors.ServerError{OrgError: "Something went wrong while operating"}
		} else if !correct {
			return &customerrors.ValidationError{OrgError: "Invalid image format"}
		}
	}
	return nil
}
