package usecases

import (
	"errors"
	customerrors "espectro/custom_errors"
	"espectro/entity"
	"espectro/pkg"
	"espectro/repository"

	"gorm.io/gorm"
)

type PartnerUsecases struct {
	repo repository.PartnerRepo
}

func NewPartnerUsecases(repo repository.PartnerRepo) PartnerUsecases {
	return PartnerUsecases{repo: repo}
}

func (p PartnerUsecases) AddPartner(partner entity.PartnerCreateEntity) (entity.PartnerEntity, error) {

	emptyEntity := entity.PartnerEntity{}
	if err := pkg.ValidateName(partner.Name); err != nil {
		return emptyEntity, &customerrors.ValidationError{OrgError: err.Error()}
	}

	newPartner, insertionErr := p.repo.AddPartner(entity.PartnerEntity{
		Name:    partner.Name,
		LogoUrl: partner.LogoUrl,
	})

	if insertionErr != nil {
		return emptyEntity, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}

	return newPartner, nil
}

func (p PartnerUsecases) UpdatePartner(newPartner entity.PartnerUpdateEntity) error {

	if !pkg.ValidateUUID(newPartner.Id) {
		return &customerrors.ValidationError{OrgError: "Invalid partner id"}
	}

	err := p.repo.UpdatePartner(newPartner)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &customerrors.NotFoundError{OrgError: "Partner does not exist"}
	} else if err != nil {
		return &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}

	return nil

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
