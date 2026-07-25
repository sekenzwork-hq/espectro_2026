package models

type SpectrumMediaModel struct {
	LogoUrl                *string
	LogoPublicId           *string
	PrevLogoPublicId       *string
	VideoUrl               *string
	VideoPublicId          *string
	PrevVideoPublicId      *string
	ImageUrls              []string
	PreviousImagePublicIds []string
	NewImagePublicIds      []string
}
