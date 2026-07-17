package entity

type EventSponsorEntity struct {
	SponsorId string `gorm:"column:sponsor_id"`
	EventId   string `gorm:"column:event_id"`
}
