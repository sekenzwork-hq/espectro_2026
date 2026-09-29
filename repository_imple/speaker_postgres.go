package repositoryimple

import (
	"espectro/entity"
	"time"

	"gorm.io/gorm"
)

type SpeakerPostgresRepo struct {
	db *gorm.DB
}

func NewSpeakerPostgresRepo(db *gorm.DB) SpeakerPostgresRepo {
	return SpeakerPostgresRepo{db: db}
}

func (s SpeakerPostgresRepo) CreateSpeaker(newSpeaker entity.SpeakerDBCreateEntity) (entity.SpeakerDBRetrieveEntity, error) {

	var createdSpeaker entity.SpeakerDBRetrieveEntity

	err := s.db.Raw(
		`
		INSERT INTO speakers(id,fullname,profile_pic_url,is_featured,bio,country,phone,email)
		VALUES(?,?,?,?,?,?,?,?)

		RETURNING id,fullname,profile_pic_url,is_featured,bio,country,phone,email,created_at;
		`,
		newSpeaker.Id,
		newSpeaker.Fullname,
		newSpeaker.ProfilePicURL,
		newSpeaker.IsFeatured,
		newSpeaker.Bio,
		newSpeaker.Country,
		newSpeaker.PhoneNumber,
		newSpeaker.Email,
	).Scan(&createdSpeaker).Error

	return createdSpeaker, err
}

func (s SpeakerPostgresRepo) UpdateSpeaker(speakerId string, newSpeaker entity.SpeakerDBUpdateEntity) (entity.SpeakerDBRetrieveEntity, error) {

	var updatedSpeaker entity.SpeakerDBRetrieveEntity

	out := s.db.Raw(
		`
		UPDATE speakers SET
		fullname=COALESCE(?,fullname),
		bio=COALESCE(?,bio),
		country=COALESCE(?,country),
		phone=COALESCE(?,phone),
		email=COALESCE(?,email),
		profile_pic_url=COALESCE(?,profile_pic_url),
		is_featured=COALESCE(?,is_featured)

		WHERE id=? AND deleted_at IS NULL

		RETURNING
		id,
		fullname,
		bio,
		country,
		phone,
		email,
		profile_pic_url,
		is_featured,
		created_at;
		`,
		newSpeaker.Fullname,
		newSpeaker.Bio,
		newSpeaker.Country,
		newSpeaker.PhoneNumber,
		newSpeaker.Email,
		newSpeaker.ProfilePicURL,
		newSpeaker.IsFeatured,
		speakerId,
	).Scan(&updatedSpeaker)

	if out.RowsAffected == 0 {
		return entity.SpeakerDBRetrieveEntity{}, gorm.ErrRecordNotFound
	}

	return updatedSpeaker, out.Error
}

func (s SpeakerPostgresRepo) DeleteSpeaker(speakerId string) error {

	out := s.db.
		Table("speakers").
		Where("id=? AND deleted_at IS NULL", speakerId).
		UpdateColumn("deleted_at", time.Now().UTC())

	if out.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return out.Error
}

func (s SpeakerPostgresRepo) SpeakerExists(speakerId string) (bool, error) {

	var exists bool

	err := s.db.Raw(`
		SELECT EXISTS( SELECT 1 FROM speakers WHERE id=? AND deleted_at IS NULL )
	`, speakerId).Scan(&exists).Error

	return exists, err
}
