package repositoryimple

import (
	"espectro/entity"
	"time"

	"gorm.io/gorm"
)

type ExhibitionPostgresRepo struct {
	db *gorm.DB
}

func NewExhibitionPostgresRepo(db *gorm.DB) ExhibitionPostgresRepo {
	return ExhibitionPostgresRepo{db: db}
}

func (e ExhibitionPostgresRepo) CreateExhibition(exhibition entity.ExhibitionDBInputEntity) (entity.ExhibitionDBRetrieveEntity, error) {

	newExhibition := entity.ExhibitionDBRetrieveEntity{}
	err := e.db.Raw(
		`INSERT INTO exhibitions
		 (
			id,
			event_id,
			category,
			organization_id,
			item_title,
			item_image_urls,
			item_description,
			status
		 ) 
		VALUES

		 (?,?,?,?,?,?,?,?)

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
		exhibition.Id,
		exhibition.EventId,
		exhibition.Category,
		exhibition.OrganizationId,
		exhibition.ItemTitle,
		exhibition.ItemImageUrls,
		exhibition.ItemDescription,
		exhibition.Status,
	).Scan(&newExhibition).Error

	return newExhibition, err
}

func (e ExhibitionPostgresRepo) UpdateExhibitionFromUserSide(exhibitionId string, newExhibition entity.ExhibitionDBInputEntity) (entity.ExhibitionDBRetrieveEntity, error) {

	updatedExhibition := entity.ExhibitionDBRetrieveEntity{}

	out := e.db.Raw(
		`UPDATE exhibitions SET 

		 event_id=COALESCE(?,event_id),
		 category=COALESCE(?,category),
		 organization_id=COALESCE(?,organization_id),
		 item_title=COALESCE(?,item_title),
		 item_description=COALESCE(?,item_description),
		 item_image_urls=COALESCE(?,item_image_urls)

		 WHERE id = ?

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

		newExhibition.EventId,
		newExhibition.Category,
		newExhibition.OrganizationId,
		newExhibition.ItemTitle,
		newExhibition.ItemDescription,
		newExhibition.ItemImageUrls,
		exhibitionId,
	).Scan(&updatedExhibition)

	if out.RowsAffected == 0 {
		return entity.ExhibitionDBRetrieveEntity{}, gorm.ErrRecordNotFound
	}

	return updatedExhibition, out.Error
}

func (e ExhibitionPostgresRepo) UpdateExhibitionFromAdminSide(exhibitionId string, newExhibition entity.ExhibitionDBInputEntity) (entity.ExhibitionDBRetrieveEntity, error) {

	updatedExhibition := entity.ExhibitionDBRetrieveEntity{}

	out := e.db.Raw(
		`UPDATE exhibitions SET 

		 token_number=COALESCE(?,token_number),
		 booth_number=COALESCE(?,booth_number),
		 available_sqft=COALESCE(?,available_sqft),
		 status=COALESCE(?,status),
		 approved_by=COALESCE(?,approved_by),
		 assigned_staff=COALESCE(?,assigned_staff)

		 WHERE id = ?

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

		newExhibition.TokenNumber,
		newExhibition.BoothNumber,
		newExhibition.AvailableSqft,
		newExhibition.Status,
		newExhibition.ApprovedBy,
		newExhibition.AssignedStaffId,
		exhibitionId,
	).Scan(&updatedExhibition)

	if out.RowsAffected == 0 {
		return entity.ExhibitionDBRetrieveEntity{}, gorm.ErrRecordNotFound
	}

	return updatedExhibition, out.Error
}

func (e ExhibitionPostgresRepo) DeleteExhibition(exhibitionId string) error {

	out := e.db.Table("exhibitions").Where("id=?", exhibitionId).Update("deleted_at", time.Now().UTC())

	if out.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
