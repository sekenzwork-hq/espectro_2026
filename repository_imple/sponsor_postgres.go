package repositoryimple

import (
	"espectro/entity"

	"gorm.io/gorm"
)

type SponsorPostgresRepo struct {
	db *gorm.DB
}

func NewSponsorPostgresRepo(db *gorm.DB) SponsorPostgresRepo {
	return SponsorPostgresRepo{db: db}
}

func (s SponsorPostgresRepo) CreateSponsor(sponsor entity.SponsorEntity) (entity.SponsorEntity, error) {

	err := s.db.
		Table("sponsors").
		Create(&sponsor).Error

	return sponsor, err
}

func (s SponsorPostgresRepo) UpdateSponsor(sponsor entity.SponsorUpdateEntity) (entity.SponsorEntity, error) {

	var entity entity.SponsorEntity

	out := s.db.Raw(
		`
		UPDATE sponsors SET
		name=COALESCE(?,name),
		amount=COALESCE(?,amount),
		profile_or_org_url=COALESCE(?,profile_or_org_url),
		sponsored_type=COALESCE(?,sponsored_type)
		
		WHERE id=? AND deleted_at IS NULL
		RETURNING id,name,amount,profile_or_org_url,sponsored_type
		`, sponsor.Name, sponsor.Amount, sponsor.ProfileOrOrgUrl, sponsor.Sponsored, sponsor.Id,
	).Scan(&entity)

	if out.Error != nil {
		return entity, out.Error
	} else if out.RowsAffected == 0 {
		return entity, gorm.ErrRecordNotFound
	}

	return entity, nil
}
