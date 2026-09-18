package entity

import (
	"espectro/enums"
	"mime/multipart"
	"time"

	"github.com/lib/pq"
)

type ExhibitionDBRetrieveEntity struct {
	Id              string                 `json:"id" gorm:"column:id;type uuid;primaryKey"`
	EventId         string                 `json:"event_id" gorm:"column:event_id"`
	TokenNumber     *int                   `json:"token_number" gorm:"column:token_number"`
	Category        string                 `json:"category" gorm:"column:category"`
	OrganizationId  string                 `json:"organization_id" gorm:"column:organization_id"`
	BoothNumber     *int                   `json:"booth_number" gorm:"column:booth_number"`
	AvailableSqft   *float32               `json:"available_sqft" gorm:"column:available_sqft"`
	AssignedStaffId *string                `json:"assigned_staff_id" gorm:"column:assigned_staff"`
	ApprovedBy      *string                `json:"approved_by" gorm:"column:approved_by"`
	ItemTitle       string                 `json:"item_title" gorm:"column:item_title"`
	ItemImageUrls   pq.StringArray         `json:"item_image_urls" gorm:"column:item_image_urls; type:text[]"`
	ItemDescription string                 `json:"item_description" gorm:"column:item_description"`
	Status          enums.ExhibitionStatus `json:"status" gorm:"column:status"`
	CreatedAt       time.Time              `json:"created_at" gorm:"column:created_at"`
}

type ExhibitionDBInputEntity struct {
	Id              *string                 `json:"id" gorm:"column:id"`
	EventId         *string                 `json:"event_id" gorm:"column:event_id"`
	TokenNumber     *int                    `json:"token_number" gorm:"column:token_number"`
	Category        *string                 `json:"category" gorm:"column:category"`
	OrganizationId  *string                 `json:"organization_id" gorm:"column:organization_id"`
	BoothNumber     *int                    `json:"booth_number" gorm:"column:booth_number"`
	AvailableSqft   *float32                `json:"available_sqft" gorm:"column:available_sqft"`
	AssignedStaffId *string                 `json:"assigned_staff_id" gorm:"column:assigned_staff"`
	ApprovedBy      *string                 `json:"approved_by" gorm:"column:approved_by"`
	ItemTitle       *string                 `json:"item_title" gorm:"column:item_title"`
	ItemImageUrls   *pq.StringArray         `json:"item_image_urls" gorm:"column:item_image_urls; type:text[]"`
	ItemDescription *string                 `json:"item_description" gorm:"column:item_description"`
	Status          *enums.ExhibitionStatus `json:"status" gorm:"column:status"`
}

type ExhibitionRawEntity struct {
	Id              *string
	EventId         string
	TokenNumber     *int
	Category        string
	OrganizationId  string
	BoothNumber     *int
	AvailableSqft   *float32
	AssignedStaffId *string
	ApprovedBy      *string
	ItemTitle       string
	ItemImages      []*multipart.FileHeader
	ItemDescription string
	Status          enums.ExhibitionStatus
}

type ExhibitionRawUpdateEntity struct {
	EventId         *string
	TokenNumber     *int
	Category        *string
	OrganizationId  *string
	BoothNumber     *int
	AvailableSqft   *float32
	AssignedStaffId *string
	ApprovedBy      *string
	ItemTitle       *string
	ItemImages      []*multipart.FileHeader
	ItemDescription *string
	Status          *enums.ExhibitionStatus
}

type ExhibitionUpdateAdminEntity struct {
	Id              string                  `json:"exhibition_id"`
	TokenNumber     *int                    `json:"token_number"`
	BoothNumber     *int                    `json:"booth_number"`
	AvailableSqft   *float32                `json:"available_sqft"`
	AssignedStaffId *string                 `json:"assigned_staff_id"`
	ApprovedBy      *string                 `json:"approved_by"`
	Status          *enums.ExhibitionStatus `json:"status"`
}
