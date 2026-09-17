package repository

import "espectro/entity"

type ExhibitionRepo interface {
	CreateExhibition(exhibition entity.ExhibitionDBCreateEntity) (entity.ExhibitionEntity, error)
}
