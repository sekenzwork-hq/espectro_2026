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

type PartnerUsecases struct {
	repo      repository.PartnerRepo
	mediaRepo repository.MediaServiceRepo
}

func NewPartnerUsecases(repo repository.PartnerRepo, mediaRepo repository.MediaServiceRepo) PartnerUsecases {
	return PartnerUsecases{repo: repo, mediaRepo: mediaRepo}
}

func (p PartnerUsecases) CreatePartner(name string, logo *multipart.FileHeader) (entity.PartnerEntity, error) {

	emptyEntity := entity.PartnerEntity{}
	if err := pkg.ValidateName(name); err != nil {
		return emptyEntity, &customerrors.ValidationError{OrgError: err.Error()}
	}

	partnerId := uuid.New().String()
	var logoUrl *string

	if logo != nil {
		url, err := p.mediaRepo.UploadFile(logo, "partner/"+partnerId)
		if err != nil {
			return emptyEntity, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
		}
		logoUrl = &url
	}

	newPartner, insertionErr := p.repo.CreatePartner(entity.PartnerEntity{
		Name:    name,
		LogoUrl: logoUrl,
		Id:      partnerId,
	})

	if insertionErr != nil {
		if logoUrl != nil {
			p.mediaRepo.DeleteFolderWithFiles("partner/", partnerId)

		}
		return emptyEntity, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}

	return newPartner, nil
}

func (p PartnerUsecases) UpdatePartner(partnerId string, name *string, logo *multipart.FileHeader) (entity.PartnerEntity, error) {

	emptyPartner := entity.PartnerEntity{}
	if !pkg.ValidateUUID(partnerId) {
		return emptyPartner, &customerrors.ValidationError{OrgError: "Invalid partner id"}
	}

	var logoUrl *string

	if logo != nil {
		url, err := p.mediaRepo.UploadFile(logo, "partner/"+partnerId)
		if err != nil {
			return emptyPartner, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
		}
		logoUrl = &url
	}

	newPartner, err := p.repo.UpdatePartner(partnerId, name, logoUrl)
	if err != nil {
		if logoUrl != nil {
			p.mediaRepo.DeleteFolderWithFiles("partner/", partnerId)
		}
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return emptyPartner, &customerrors.NotFoundError{OrgError: "Partner does not exist"}
	} else if err != nil {
		return emptyPartner, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}

	return newPartner, nil

}

func (p PartnerUsecases) DeletePartner(partnerId string) error {

	if !pkg.ValidateUUID(partnerId) {
		return &customerrors.ValidationError{OrgError: "Invalid partner id"}
	}

	err := p.repo.DeletePartner(partnerId)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &customerrors.NotFoundError{OrgError: "Partner does not exist"}
	} else if err != nil {
		return &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}

	return nil
}

func (p PartnerUsecases) RetrievePartner(limit int, page int) ([]entity.PartnerEntity, error) {

	offset := pkg.GetOffset(limit, page)
	partners, err := p.repo.RetrievePartner(limit, offset)

	if err != nil {
		return partners, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}

	return partners, nil
}
