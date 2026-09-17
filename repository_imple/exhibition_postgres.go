package repositoryimple

import (
	"espectro/entity"

	"gorm.io/gorm"
)

type ExhibitionPostgresRepo struct {
	db *gorm.DB
}

func NewExhibitionPostgresRepo(db *gorm.DB) ExhibitionPostgresRepo {
	return ExhibitionPostgresRepo{db: db}
}

func (e ExhibitionPostgresRepo) CreateExhibition(exhibition entity.ExhibitionDBCreateEntity) (entity.ExhibitionEntity, error) {

	newExhibitions := entity.ExhibitionEntity{}
	err := e.db.Table("exhibitions").Raw(
		`INSERT INTO exhibitions
		 (
			event_id,
			category,
			organization_id,
			item_title,
			item_image_urls,
			item_description,
			status
		 ) 
		VALUES

		 (
			?,
			?,
			?,
			?,
			?,
			?,
			?
		 )

		RETURNING 
			id,
			event_id,
			token_number,
			category,
			organization_id,
			booth_number,
			available_sqft,
			assigned_staff,
			approved_by,
			item_title,
			item_image_urls,
			item_description,
			status,
			created_at;
		`,
		exhibition.EventId,
		exhibition.Category,
		exhibition.OrganizationId,
		exhibition.ItemTitle,
		exhibition.ItemImageUrls,
		exhibition.ItemDescription,
		exhibition.Status,
	).Scan(&newExhibitions).Error

	return newExhibitions, err
}
