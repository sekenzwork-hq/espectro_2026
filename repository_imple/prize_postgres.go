package repositoryimple

import (
	"espectro/entity"

	"gorm.io/gorm"
)

type PrizePostgresRepo struct {
	db *gorm.DB
}

func NewPrizePostgresRepo(db *gorm.DB) PrizePostgresRepo {
	return PrizePostgresRepo{db: db}
}

func (p PrizePostgresRepo) CreatePrize(prize entity.PrizeDBCreateEntity) (entity.PrizeDBRetrieveEntity, error) {

	var createdPrize entity.PrizeDBRetrieveEntity

	err := p.db.Raw(`
	
	INSERT INTO prize(id,title,description,amount,logo_url,event_id,type)

	VALUES(?,?,?,?,?,?,?)

	RETURNING 
	id,
	title,
	description,
	amount,
	logo_url,
	type,
	created_at;
	
	`,
		prize.Id,
		prize.Title,
		prize.Description,
		prize.Amount,
		prize.LogoURL,
		prize.EventId,
		prize.Type,
	).Scan(&createdPrize).Error

	return createdPrize, err
}

func (p PrizePostgresRepo) UpdatePrize(prizeId string, newPrize entity.PrizeDBUpdateEntity) (entity.PrizeDBRetrieveEntity, error) {

	var updatedPrize entity.PrizeDBRetrieveEntity

	out := p.db.Raw(
		`
		UPDATE prize SET

		title=COALESCE(?,title),
		description=COALESCE(?,description),
		logo_url=COALESCE(?,logo_url),
		amount=COALESCE(?,amount),
		event_id=COALESCE(?,event_id),
		type=COALESCE(?,type)

		WHERE id=? AND deleted_at IS NULL

		RETURNING
		id,
		title,
		description,
		logo_url,
		amount,
		event_id,
		type,
		created_at;

		`,
		newPrize.Title,
		newPrize.Description,
		newPrize.LogoURL,
		newPrize.Amount,
		newPrize.EventId,
		newPrize.Type,
		prizeId,
	).Scan(&updatedPrize)

	if out.RowsAffected == 0 {
		return updatedPrize, gorm.ErrRecordNotFound
	}

	return updatedPrize, out.Error
}
