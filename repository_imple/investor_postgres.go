package repositoryimple

import (
	"espectro/entity"
	"time"

	"gorm.io/gorm"
)

type InvestorPostgresRepo struct {
	db *gorm.DB
}

func NewInvestorPostgresRepo(db *gorm.DB) InvestorPostgresRepo {
	return InvestorPostgresRepo{db: db}
}

func (i InvestorPostgresRepo) CreateInvestor(investor entity.InvestorEntity) (entity.InvestorEntity, error) {

	err := i.db.
		Table("investors").
		Create(&investor).Error

	return investor, err
}

func (i InvestorPostgresRepo) UpdateInvestor(investorId string, newInvestor entity.InvestorUpdateEntity) (entity.InvestorEntity, error) {

	var investor entity.InvestorEntity

	out := i.db.
		Raw(`UPDATE investors SET 
		name=COALESCE(?,name),
		website_url=COALESCE(?,website_url),
		phone=COALESCE(?,phone),
		email=COALESCE(?,email),
		logo_url=COALESCE(?,logo_url)

		WHERE id=? AND deleted_at IS NULL
		RETURNING id,name,website_url,phone,email,logo_url,created_at
		`,
			newInvestor.Name, newInvestor.WebsiteUrl, newInvestor.PhoneNumber, newInvestor.Email, newInvestor.LogoUrl, investorId,
		).
		Scan(&investor)

	if out.Error != nil {
		return investor, out.Error
	} else if out.RowsAffected == 0 {
		return investor, gorm.ErrRecordNotFound
	}

	return investor, nil

}

func (i InvestorPostgresRepo) DeleteInvestor(investorId string) error {

	out := i.db.
		Table("investors").
		Where("id=? AND deleted_at IS NULL", investorId).
		Update("deleted_at", time.Now().UTC())

	if out.Error != nil {
		return out.Error
	} else if out.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (i InvestorPostgresRepo) RetrieveInvestors(limit int, offset int) ([]entity.InvestorEntity, error) {

	var investors []entity.InvestorEntity

	err := i.db.
		Table("investors").
		Where("deleted_at IS NULL").
		Offset(offset).
		Limit(limit).
		Scan(&investors).Error

	return investors, err
}
