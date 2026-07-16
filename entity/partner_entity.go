package entity

type PartnerEntity struct {
	Id        string  `json:"id" gorm:"primaryKey;type uuid;default:gen_random_uuid()"`
	Name      string  `json:"name" gorm:"column:name"`
	LogoUrl   *string `json:"logo_url" gorm:"column:logo_url"`
	CreatedAt string  `json:"created_at" gorm:"type timestampz;default:now()"`
}

type PartnerCreateEntity struct {
	Name    string  `json:"name"`
	LogoUrl *string `json:"logo_url"`
}

type PartnerUpdateEntity struct {
	Id      string  `json:"id"`
	Name    *string `json:"name" gorm:"column:name"`
	LogoUrl *string `json:"logo_url" gorm:"column:logo_url"`
}
