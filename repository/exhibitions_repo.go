package repository

import "espectro/entity"

type ExhibitionRepo interface {
	CreateExhibition(exhibition entity.ExhibitionDBInputEntity) (entity.ExhibitionDBRetrieveEntity, error)
	UpdateExhibitionFromUserSide(exhibitionId string, newExhibition entity.ExhibitionDBInputEntity) (entity.ExhibitionDBRetrieveEntity, error)
	UpdateExhibitionFromAdminSide(exhibitionId string, newExhibition entity.ExhibitionDBInputEntity) (entity.ExhibitionDBRetrieveEntity, error)
	DeleteExhibition(exhibitionId string) error
}
