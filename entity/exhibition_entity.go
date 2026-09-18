package entity

import (
	"espectro/enums"
	"mime/multipart"
	"time"

	"github.com/lib/pq"
)

type ExhibitionEntity struct {
	Id              string                 `json:"id" gorm:"column:id;type uuid;primaryKey"`
	EventId         string                 `json:"event_id" gorm:"column:event_id"`
	TokenNumber     *int                   `json:"token_number" gorm:"column:token_number"`
	Category        string                 `json:"category" gorm:"column:category"`
	OrganizationId  string                 `json:"organization_id" gorm:"column:organization_id"`
	BoothNumber     *int                   `json:"booth_number" gorm:"column:booth_number"`
	AvailableSqft   *float64               `json:"available_sqft" gorm:"column:available_sqft"`
	AssignedStaffId *string                `json:"assigned_staff_id" gorm:"column:assigned_staff"`
	ApprovedBy      *string                `json:"approved_by" gorm:"column:approved_by"`
	ItemTitle       string                 `json:"item_title" gorm:"column:item_title"`
	ItemImageUrls   pq.StringArray         `json:"item_image_urls" gorm:"column:item_image_urls; type:text[]"`
	ItemDescription string                 `json:"item_description" gorm:"column:item_description"`
	Status          enums.ExhibitionStatus `json:"status" gorm:"column:status"`
	CreatedAt       time.Time              `json:"created_at" gorm:"column:created_at"`
}

type ExhibitionCreateEntity struct {
	EventId         string
	Category        string
	OrganizationId  string
	ItemTitle       string
	ItemImages      []*multipart.FileHeader
	ItemDescription string
}

type ExhibitionDBCreateEntity struct {
	Id              string                 `gorm:"column:id;type uuid;primaryKey"`
	EventId         string                 `gorm:"column:event_id"`
	Category        string                 `gorm:"column:category"`
	OrganizationId  string                 `gorm:"column:organization_id"`
	ItemTitle       string                 `gorm:"column:item_title"`
	ItemImageUrls   pq.StringArray         `gorm:"column:item_image_urls; type:text[]"`
	ItemDescription string                 `gorm:"column:item_description"`
	Status          enums.ExhibitionStatus `gorm:"column:status"`
}
