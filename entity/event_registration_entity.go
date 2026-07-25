package entity

import (
	"espectro/enums"
	"time"
)

type EventRegistrationEntity struct {
	Id        string               `json:"id" gorm:"column:id;type uuid;default:gen_random_uuid();primaryKey"`
	UserId    string               `json:"user_id" gorm:"column:user_id"`
	EventId   string               `json:"event_id" gorm:"column:event_id"`
	Status    enums.EventRegStatus `json:"registration_status" gorm:"column:status"`
	CheckIn   *time.Time           `json:"check_in" gorm:"column:check_in"`
	CheckOut  *time.Time           `json:"check_out" gorm:"column:check_out"`
	CreatedAt time.Time            `json:"created_at" gorm:"column:created_at;type timestampz;default:now()"`
}

type EventRegistrationFromJsonEntity struct {
	UserId  string `json:"user_id"`
	EventId string `json:"event_id"`
}

type ChangeRegistrationStatusEntity struct {
	Status enums.EventRegStatus `json:"status"`
	Id     string               `json:"registration_id"`
}
