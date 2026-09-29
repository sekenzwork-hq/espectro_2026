package entity

import "mime/multipart"

type SpeakerDBCreateEntity struct {
	Id            string  `json:"id" gorm:"column:id; type:uuid()"`
	Fullname      string  `json:"fullname" gorm:"column:fullname"`
	ProfilePicURL *string `json:"profile_pic_url" gorm:"column:profile_pic_url"`
	IsFeatured    bool    `json:"is_featured" gorm:"column:is_featured"`
	Bio           string  `json:"bio" gorm:"column:bio"`
	Country       string  `json:"country" gorm:"column:country"`
	PhoneNumber   string  `json:"phone_number" gorm:"column:phone"`
	Email         string  `json:"email" gorm:"column:email"`
}

type SpeakerDBRetrieveEntity struct {
	Id            string  `json:"id" gorm:"column:id; type:uuid()"`
	Fullname      string  `json:"fullname" gorm:"column:fullname"`
	ProfilePicURL *string `json:"profile_pic_url" gorm:"column:profile_pic_url"`
	IsFeatured    bool    `json:"is_featured" gorm:"column:is_featured"`
	Bio           string  `json:"bio" gorm:"column:bio"`
	Country       string  `json:"country" gorm:"column:country"`
	PhoneNumber   string  `json:"phone_number" gorm:"column:phone"`
	Email         string  `json:"email" gorm:"column:email"`
	CreatedAt     string  `json:"created_at" gorm:"column:created_at"`
}

type SpeakerDBUpdateEntity struct {
	Fullname      *string `json:"fullname" gorm:"column:fullname"`
	ProfilePicURL *string `json:"profile_pic_url" gorm:"column:profile_pic_url"`
	IsFeatured    *bool   `json:"is_featured" gorm:"column:is_featured"`
	Bio           *string `json:"bio" gorm:"column:bio"`
	Country       *string `json:"country" gorm:"column:country"`
	PhoneNumber   *string `json:"phone_number" gorm:"column:phone"`
	Email         *string `json:"email" gorm:"column:email"`
}

type SpeakerCreateEntity struct {
	Fullname    string
	ProfilePic  *multipart.FileHeader
	IsFeatured  bool
	Bio         string
	Country     string
	PhoneNumber string
	Email       string
}
type SpeakerUpdateEntity struct {
	Fullname    *string
	ProfilePic  *multipart.FileHeader
	IsFeatured  *bool
	Bio         *string
	Country     *string
	PhoneNumber *string
	Email       *string
}
