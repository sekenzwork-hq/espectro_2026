package repository

import "espectro/entity"

type InvestorRepo interface {
	AddInvestor(investor entity.InvestorEntity) (entity.InvestorEntity, error)
}
