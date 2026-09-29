package repository

import "espectro/entity"

type PrizeRepo interface {
	CreatePrize(prize entity.PrizeDBCreateEntity) (entity.PrizeDBRetrieveEntity, error)
	UpdatePrize(prizeId string, newPrize entity.PrizeDBUpdateEntity) (entity.PrizeDBRetrieveEntity, error)
	DeletePrize(prizeId string) error
	RetrievePrize(prizeId string) (entity.PrizeDBRetrieveEntity, error)
}
