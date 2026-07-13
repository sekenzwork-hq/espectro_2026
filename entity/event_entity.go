package entity

type EventEntity struct {
	Id          string `json:"id" gorm:"type uuid;default:gen_random_uuid();primaryKey"`
	Name        string `json:"name" gorm:"column:name"`
	Description string `json:"description" gorm:"column:description"`
	SpectrumId  string `json:"spectrum_id" gorm:"column:spectrum_id"`
	Status      string `json:"status" gorm:"column:status"`
}
