package repositoryimple

import (
	"espectro/entity"
	"espectro/enums"
	"time"

	"gorm.io/gorm"
)

type CompetitionPostgresRepo struct {
	db *gorm.DB
}

func NewCompetitionPostgresRepo(db *gorm.DB) CompetitionPostgresRepo {
	return CompetitionPostgresRepo{db: db}
}

func (c CompetitionPostgresRepo) CreateCompetition(newCompetition entity.CompetitionDBCreateEntity) (entity.CompetitionDBRetrieveEntity, error) {

	var createdCompetition entity.CompetitionDBRetrieveEntity

	var rowStruct struct {
		Id               string                      `gorm:"column:id"`
		Title            string                      `gorm:"column:title"`
		Description      string                      `gorm:"column:description"`
		Rules            string                      `gorm:"column:rules"`
		Status           enums.CompetitionStatusEnum `gorm:"column:status"`
		OpeningTime      time.Time                   `gorm:"column:opening_time"`
		ClosingTime      time.Time                   `gorm:"column:closing_time"`
		ResultTime       time.Time                   `gorm:"column:result_time"`
		RegistrationFee  float32                     `gorm:"column:registration_fee"`
		LogoURL          *string                     `gorm:"column:logo_url"`
		CreatedAt        time.Time                   `gorm:"column:created_at"`
		PrizeId          string                      `gorm:"column:prize_id"`
		PrizeTitle       string                      `gorm:"column:prize_title"`
		PrizeDescription string                      `gorm:"column:prize_description"`
		PrizeAmount      float32                     `gorm:"column:prize_amount"`
		PrizeLogoURL     *string                     `gorm:"column:prize_logo_url"`
		PrizeType        enums.PrizeTypeEnum         `gorm:"column:prize_type"`
		PrizeCreatedAt   time.Time                   `gorm:"column:prize_created_at"`
		VenueId          string                      `gorm:"column:venue_id"`
		VenueCountry     string                      `gorm:"column:venue_country"`
		VenueState       string                      `gorm:"column:venue_state"`
		VenueCity        string                      `gorm:"column:venue_city"`
		VenueCreatedAt   time.Time                   `gorm:"column:venue_created_at"`
	}

	err := c.db.Raw(
		`
		WITH inserted AS(
		 	INSERT INTO competitions (
				id,
				event_id,
				title,
				description,
				rules,
				status,
				opening_time,
				closing_time,
				result_time,
				registration_fee,
				prize_id,
				venue_id,
				logo_url
			)
			
			VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)
			RETURNING *
		)
			
		SELECT i.id,
		i.title,
		i.description,
		i.rules,
		i.status,
		i.opening_time,
		i.closing_time,
		i.result_time,
		i.registration_fee,
		i.logo_url,
		p.id as prize_id,
		p.title as prize_title,
		p.description as prize_description,
		p.amount as prize_amount,
		p.logo_url as prize_logo_url,
		p.type as prize_type,
		p.created_at as prize_created_at,
		v.id as venue_id,
		v.country as venue_country,
		v.state as venue_state,
		v.city as venue_city,
		v.created_at as venue_created_at


		FROM inserted i
		LEFT JOIN prize p on p.id=i.prize_id
		LEFT JOIN venue v on v.id=i.venue_id
		`,
		newCompetition.Id,
		newCompetition.EventId,
		newCompetition.Title,
		newCompetition.Description,
		newCompetition.Rules,
		newCompetition.Status,
		newCompetition.OpeningTime,
		newCompetition.ClosingTime,
		newCompetition.ResultTime,
		newCompetition.RegistrationFee,
		newCompetition.PrizeId,
		newCompetition.VenueId,
		newCompetition.LogoURL,
	).Scan(&rowStruct).Error

	createdCompetition = entity.CompetitionDBRetrieveEntity{
		Id:              rowStruct.Id,
		Title:           rowStruct.Title,
		Description:     rowStruct.Description,
		Rules:           rowStruct.Rules,
		Status:          rowStruct.Status,
		OpeningTime:     rowStruct.OpeningTime,
		ClosingTime:     rowStruct.ClosingTime,
		ResultTime:      rowStruct.ResultTime,
		RegistrationFee: rowStruct.RegistrationFee,
		LogoURL:         rowStruct.LogoURL,
		CreatedAt:       rowStruct.CreatedAt,
		Prize: entity.PrizeDBRetrieveEntity{
			Id:          rowStruct.PrizeId,
			Title:       rowStruct.PrizeTitle,
			Description: rowStruct.Description,
			Amount:      rowStruct.PrizeAmount,
			LogoURL:     rowStruct.PrizeLogoURL,
			Type:        rowStruct.PrizeType,
			CreatedAt:   rowStruct.PrizeCreatedAt,
		},
		Venue: entity.VenueEntity{
			Id:        rowStruct.VenueId,
			Country:   rowStruct.VenueCountry,
			State:     rowStruct.VenueState,
			City:      rowStruct.VenueCity,
			CreatedAt: rowStruct.CreatedAt,
		},
	}

	return createdCompetition, err
}
