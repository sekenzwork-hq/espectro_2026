package repositoryimple

import (
	"espectro/entity"
	"espectro/enums"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type SpectrumPostgresRepo struct {
	db *gorm.DB
}

func NewSpectrumPostgresRepo(db *gorm.DB) SpectrumPostgresRepo {
	return SpectrumPostgresRepo{db: db}
}

func (s SpectrumPostgresRepo) CreateSpectrum(newSpectrum entity.SpectrumEntity) (entity.SpectrumEntity, error) {
	out := s.db.Table("spectrums").Create(&newSpectrum)
	return newSpectrum, out.Error
}

func (s SpectrumPostgresRepo) UpdateSpectrum(
	spectrumId string,
	name *string,
	shortDescription *string,
	description *string,
	status *enums.SpectrumStatus,
	logoUrl *string,
	videoUrl *string,
	imageUrls []string,
) (entity.SpectrumEntity, error) {

	var newSpectrum entity.SpectrumEntity
	var pgImageUrls pq.StringArray

	if len(imageUrls) != 0 {
		pgImageUrls = pq.StringArray(imageUrls)
	}

	out := s.db.
		Table("spectrums").
		Raw(`UPDATE spectrums SET 
		name=COALESCE(?,name),
		short_description=COALESCE(?,short_description),
		description=COALESCE(?,description),
		status=COALESCE(?,status),
		logo_url=COALESCE(?,logo_url),
		video_url=COALESCE(?,video_url),
		image_urls=COALESCE(?,image_urls) 

		WHERE id=? AND deleted_at IS NULL

		RETURNING id,name,short_description,description,status,logo_url,video_url,image_urls,created_at
		`,
			name, shortDescription, description, status, logoUrl, videoUrl, pgImageUrls, spectrumId,
		).
		Scan(&newSpectrum)

	if out.Error != nil {
		return newSpectrum, out.Error
	} else if out.RowsAffected == 0 {
		return newSpectrum, gorm.ErrRecordNotFound
	}

	return newSpectrum, nil
}

func (s SpectrumPostgresRepo) DeleteSpectrum(spectrumId string) error {

	out := s.db.
		Table("spectrums").
		Where("id=? AND deleted_at IS NULL", spectrumId).
		Update("deleted_at", time.Now().UTC())

	if out.Error != nil {
		return out.Error
	} else if out.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (s SpectrumPostgresRepo) RetrieveSpectrums(offset int, limit int) ([]entity.SpectrumEntity, error) {

	var spectrums []entity.SpectrumEntity

	out := s.db.
		Table("spectrums").
		Select("id,name,short_description,description,status,total_events,logo_url,video_url,image_urls,created_at").
		Where("deleted_at IS NULL").
		Offset(offset).
		Limit(limit).
		Scan(&spectrums)

	return spectrums, out.Error
}

func (s SpectrumPostgresRepo) CheckSpectrumExists(spectrumId string) (bool, error) {

	var exists bool
	err := s.db.
		Raw("SELECT EXISTS (SELECT 1 FROM spectrums WHERE id=? AND deleted_at IS NULL)", spectrumId).
		Scan(&exists).Error

	return exists, err
}

func (s SpectrumPostgresRepo) IncrementTotalEventsCountBy1(spectrumId string) error {

	out := s.db.
		Table("spectrums").
		Where("id=? AND deleted_at IS NULL", spectrumId).
		UpdateColumn("total_events", gorm.Expr("total_events+1"))

	if out.Error != nil {
		return out.Error
	} else if out.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (s SpectrumPostgresRepo) DecrementTotalEventsCountBy1(spectrumId string) error {
	out := s.db.
		Table("spectrums").
		Where("id=? AND deleted_at IS NULL", spectrumId).
		UpdateColumn("total_events", gorm.Expr(`
			CASE 
		      WHEN total_events != 0 THEN total_events-1
			END
		`))

	if out.Error != nil {
		return out.Error
	}

	return nil
}
