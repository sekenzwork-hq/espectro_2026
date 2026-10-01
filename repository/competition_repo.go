package repository

import "espectro/entity"

type CompetitionRepo interface {
	CreateCompetition(newCompetition entity.CompetitionDBCreateEntity) (entity.CompetitionDBRetrieveEntity, error)
	UpdateCompetition(competitionId string, newCompetition entity.CompetitionDBUpdateEntity) (entity.CompetitionDBRetrieveEntity, error)
	DeleteCompetition(competitionId string) error
	RetrieveCompetitions(eventId string, offset int) ([]entity.CompetitionDBRetrieveEntity, error)
}
