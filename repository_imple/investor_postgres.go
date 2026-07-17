package repositoryimple

import (
	"espectro/entity"

	"gorm.io/gorm"
)

type InvestorPostgresRepo struct {
	db *gorm.DB
}

func NewInvestorPostgresRepo(db *gorm.DB) InvestorPostgresRepo {
	return InvestorPostgresRepo{db: db}
}

func (i InvestorPostgresRepo) AddInvestor(investor entity.InvestorEntity) (entity.InvestorEntity, error) {

	err := i.db.
		Table("investors").
		Create(&investor).Error

	return investor, err
}
