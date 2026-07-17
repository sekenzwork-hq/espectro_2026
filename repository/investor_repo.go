package repository

import "espectro/entity"

type InvestorRepo interface {
	AddInvestor(investor entity.InvestorEntity) (entity.InvestorEntity, error)
	UpdateInvestor(investorId string, newInvestor entity.InvestorUpdateEntity) (entity.InvestorEntity, error)
	DeleteInvestor(investorId string) error
	RetrieveInvestors(limit int, offset int) ([]entity.InvestorEntity, error)
}
