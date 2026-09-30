package repository

import "espectro/entity"

type CompetitionRepo interface {
	CreateCompetition(newCompetition entity.CompetitionDBCreateEntity) (entity.CompetitionDBRetrieveEntity, error)
}
