package repositoryimple

import (
	"espectro/entity"
	"espectro/enums"

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
	imageUrls *[]string,
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
		data["image_urls"] = pq.StringArray(*(imageUrls))
	}

	out := s.db.Table("spectrums").Where("id=?", spectrumId).Updates(data)

	if out.Error != nil {
		return out.Error
	} else if out.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
