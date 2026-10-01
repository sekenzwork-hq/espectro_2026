package repositoryimple

import (
	"espectro/entity"
	"espectro/enums"
	"time"

	"gorm.io/gorm"
)

type rowStruct struct {
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

type CompetitionPostgresRepo struct {
	db *gorm.DB
}

func NewCompetitionPostgresRepo(db *gorm.DB) CompetitionPostgresRepo {
	return CompetitionPostgresRepo{db: db}
}

func (c CompetitionPostgresRepo) CreateCompetition(newCompetition entity.CompetitionDBCreateEntity) (entity.CompetitionDBRetrieveEntity, error) {

	var createdCompetition entity.CompetitionDBRetrieveEntity

	var row rowStruct
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
	).Scan(&row).Error

	createdCompetition = c.parseRowStructToCompetition(row)
	return createdCompetition, err
}

func (c CompetitionPostgresRepo) UpdateCompetition(competitionId string, newCompetition entity.CompetitionDBUpdateEntity) (entity.CompetitionDBRetrieveEntity, error) {

	var row rowStruct

	out := c.db.Raw(
		`
		WITH updated AS(
		 	UPDATE competitions SET
				event_id=COALESCE(?, event_id),
				title=COALESCE(?, title),
				description=COALESCE(?, description),
				rules=COALESCE(?, rules),
				status=COALESCE(?, status),
				opening_time=COALESCE(?, opening_time),
				closing_time=COALESCE(?, closing_time),
				result_time=COALESCE(?, result_time),
				registration_fee=COALESCE(?, registration_fee),
				prize_id=COALESCE(?, prize_id),
				venue_id=COALESCE(?, venue_id),
				logo_url=COALESCE(?, logo_url)

			WHERE id=? AND deleted_at IS NULL
			RETURNING *
		)
			
		SELECT u.id,
		u.title,
		u.description,
		u.rules,
		u.status,
		u.opening_time,
		u.closing_time,
		u.result_time,
		u.registration_fee,
		u.logo_url,
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

		FROM updated u
		LEFT JOIN prize p on p.id=u.prize_id
		LEFT JOIN venue v on v.id=u.venue_id		
		`,
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
		competitionId,
	).Scan(&row)

	if out.RowsAffected == 0 {
		return entity.CompetitionDBRetrieveEntity{}, gorm.ErrRecordNotFound
	}

	updatedCompetition := c.parseRowStructToCompetition(row)

	return updatedCompetition, out.Error
}

func (c CompetitionPostgresRepo) DeleteCompetition(competitionId string) error {

	out := c.db.
		Table("competitions").
		Where("id=? AND deleted_at IS NULL", competitionId).
		UpdateColumn("deleted_at", time.Now().UTC())

	if out.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return out.Error
}

func (c CompetitionPostgresRepo) RetrieveCompetitions(eventId string, offset int) ([]entity.CompetitionDBRetrieveEntity, error) {

	var rows []rowStruct

	var competitions []entity.CompetitionDBRetrieveEntity

	err := c.db.Raw(
		`
		SELECT 
		c.id,
		c.title,
		c.description,
		c.rules,
		c.status,
		c.opening_time,
		c.closing_time,
		c.result_time,
		c.registration_fee,
		c.logo_url,
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

		FROM competitions c
		LEFT JOIN prize p on p.id=c.prize_id
		LEFT JOIN venue v on v.id=c.venue_id
		WHERE c.event_id=? AND c.deleted_at IS NULL
		OFFSET ?
		LIMIT 50
		`,
		eventId,
		offset,
	).Scan(&rows).Limit(50).Error

	for i := range rows {
		row := rows[i]
		competition := c.parseRowStructToCompetition(row)
		competitions = append(competitions, competition)
	}

	return competitions, err
}

func (c CompetitionPostgresRepo) parseRowStructToCompetition(row rowStruct) entity.CompetitionDBRetrieveEntity {

	return entity.CompetitionDBRetrieveEntity{
		Id:              row.Id,
		Title:           row.Title,
		Description:     row.Description,
		Rules:           row.Rules,
		Status:          row.Status,
		OpeningTime:     row.OpeningTime,
		ClosingTime:     row.ClosingTime,
		ResultTime:      row.ResultTime,
		RegistrationFee: row.RegistrationFee,
		LogoURL:         row.LogoURL,
		CreatedAt:       row.CreatedAt,
		Prize: entity.PrizeDBRetrieveEntity{
			Id:          row.PrizeId,
			Title:       row.PrizeTitle,
			Description: row.Description,
			Amount:      row.PrizeAmount,
			LogoURL:     row.PrizeLogoURL,
			Type:        row.PrizeType,
			CreatedAt:   row.PrizeCreatedAt,
		},
		Venue: entity.VenueEntity{
			Id:        row.VenueId,
			Country:   row.VenueCountry,
			State:     row.VenueState,
			City:      row.VenueCity,
			CreatedAt: row.CreatedAt,
		},
	}

}
