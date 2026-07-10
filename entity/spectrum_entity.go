package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Spectrum struct {
	Id               uuid.UUID      `json:"id" gorm:"type uuid;default:gen_random_uuid();primaryKey"`
	Name             string         `json:"name" gorm:"column:name"`
	ShortDescription string         `json:"short_description" gorm:"column:short_description"`
	Description      string         `json:"description" gorm:"column:description"`
	LogoUrl          *string        `json:"logo_url" gorm:"column:logo_url"`
	ImageUrls        pq.StringArray `json:"image_urls" gorm:"type:text[]"`
	VideoUrl         *string        `json:"video_url" gorm:"column:video_url"`
	TotalEvents      int            `json:"total_events" gorm:"column:total_events"`
	CreatedAt        *time.Time     `json:"created_at" gorm:"type:timestampz;default:now()"`
	DeletedAt        *time.Time     `json:"deleted_at" gorm:"type:timestampz"`
}
