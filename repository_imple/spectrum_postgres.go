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
) error {

	data := map[string]any{}

	if name != nil {
		data["name"] = *(name)
	}
	if shortDescription != nil {
		data["short_description"] = *(shortDescription)
	}
	if description != nil {
		data["description"] = *(description)
	}
	if status != nil {
		data["status"] = *(status)
	}
	if logoUrl != nil {
		data["logo_url"] = *(logoUrl)
	}
	if videoUrl != nil {
		data["video_url"] = *(videoUrl)
	}
	if imageUrls != nil {
		data["image_urls"] = pq.StringArray(imageUrls)
	}

	out := s.db.Table("spectrums").Where("id=? AND deleted_at IS NULL", spectrumId).Updates(data)

	if out.Error != nil {
		return out.Error
	} else if out.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
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
		Raw("SELECT EXISTS (SELECT 1 FROM spectrums WHERE id=?)", spectrumId).
		Scan(&exists).Error

	return exists, err
}
