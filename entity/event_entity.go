package entity

import (
	"espectro/enums"
	"time"
)

type EventEntity struct {
	Id               string            `json:"id" gorm:"type uuid;default:gen_random_uuid();primaryKey"`
	Name             string            `json:"name" gorm:"column:name"`
	Description      string            `json:"description" gorm:"column:description"`
	SpectrumId       string            `json:"spectrum_id" gorm:"column:spectrum_id"`
	Status           enums.EventStatus `json:"status" gorm:"column:status"`
	ParticipantLimit *int              `json:"participant_limit" gorm:"column:participant_limit"`
	StartDate        *time.Time        `json:"start_date" gorm:"column:start_date"`
	EndDate          *time.Time        `json:"end_date" gorm:"column:end_date"`
	EventMode        enums.EventMode   `json:"event_mode" gorm:"column:event_mode"`
	EventType        enums.EventType   `json:"event_type" gorm:"column:event_type"`
	IsFeatured       bool              `json:"is_featured" gorm:"column:is_featured"`
	ContactEmail     string            `json:"contact_email" gorm:"contact_email"`
	VenueId          string            `json:"venue_id" gorm:"column:venue_id"`
	CreatedAt        *time.Time        `json:"created_at" gorm:"type timestampz;default:now()"`
}

type EventFromJsonEntity struct {
	Name             string            `json:"name"`
	Description      string            `json:"description"`
	SpectrumId       string            `json:"spectrum_id"`
	Status           enums.EventStatus `json:"status"`
	ParticipantLimit *int              `json:"participant_limit"`
	StartDate        *string           `json:"start_date"`
	EndDate          *string           `json:"end_date"`
	EventMode        enums.EventMode   `json:"event_mode"`
	EventType        enums.EventType   `json:"event_type"`
	IsFeatured       bool              `json:"is_featured"`
	ContactEmail     string            `json:"contact_email"`
	VenueId          string            `json:"venue_id"`
}

type EventUpdateEntity struct {
	Id               string            `json:"event_id" gorm:"-"`
	Name             string            `json:"name" gorm:"column:name"`
	Description      string            `json:"description"  gorm:"column:description"`
	SpectrumId       string            `json:"spectrum_id" gorm:"column:spectrum_id"`
	Status           enums.EventStatus `json:"status" gorm:"column:status"`
	ParticipantLimit *int              `json:"participant_limit" gorm:"column:participant_limit"`
	StartDate        *string           `json:"start_date" gorm:"column:start_date"`
	EndDate          *string           `json:"end_date" gorm:"column:end_date"`
	EventMode        enums.EventMode   `json:"event_mode" gorm:"column:event_mode"`
	EventType        enums.EventType   `json:"event_type" gorm:"column:event_type"`
	IsFeatured       bool              `json:"is_featured" gorm:"column:is_featured"`
	ContactEmail     string            `json:"contact_email" gorm:"column:contact_email"`
	VenueId          string            `json:"venue_id" gorm:"column:venue_id"`
}
