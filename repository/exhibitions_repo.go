package repository

import "espectro/entity"

type ExhibitionRepo interface {
	CreateExhibition(exhibition entity.ExhibitionDBInputEntity) (entity.ExhibitionDBRetrieveEntityFromUserSide, error)
	UpdateExhibitionFromUserSide(exhibitionId string, newExhibition entity.ExhibitionDBInputEntity) (entity.ExhibitionDBRetrieveEntityFromUserSide, error)
	UpdateExhibitionFromAdminSide(exhibitionId string, newExhibition entity.ExhibitionDBInputEntity) (entity.ExhibitionDBRetrieveEntityFromAdminSide, error)
	DeleteExhibition(exhibitionId string) error
	RetrieveExhibitionFromUserSide(userId string, offset int) ([]entity.ExhibitionDBRetrieveEntityFromUserSide, error)
	RetrieveExhibitionFromAdminSide(offset int) ([]entity.ExhibitionDBRetrieveEntityFromAdminSide, error)
}
