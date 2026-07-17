package repository

import "espectro/entity"

type InvestorRepo interface {
	AddInvestor(investor entity.InvestorEntity) (entity.InvestorEntity, error)
	UpdateInvestor(investorId string, newInvestor entity.InvestorUpdateEntity) (entity.InvestorEntity, error)
	DeleteInvestor(investorId string) error
}
