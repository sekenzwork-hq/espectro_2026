package usecases

import (
	customerrors "espectro/custom_errors"
	"espectro/entity"
	"espectro/pkg"
	"espectro/repository"
)

type PartnerUsecases struct {
	repo repository.PartnerRepo
}

func NewPartnerUsecases(repo repository.PartnerRepo) PartnerUsecases {
	return PartnerUsecases{repo: repo}
}

func (p PartnerUsecases) AddPartner(partner entity.PartnerFromJsonEntity) (entity.PartnerEntity, error) {

	emptyEntity := entity.PartnerEntity{}
	if err := pkg.ValidateName(partner.Name); err != nil {
		return emptyEntity, &customerrors.ValidationError{OrgError: err.Error()}
	}

	newPartner, insertionErr := p.repo.AddPartner(entity.PartnerEntity{
		Amount:  partner.Amount,
		Name:    partner.Name,
		LogoUrl: partner.LogoUrl,
	})

	if insertionErr != nil {
		return emptyEntity, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}

	return newPartner, nil
}
