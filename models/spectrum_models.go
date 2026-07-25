package models

type SpectrumMediaModel struct {
	LogoUrl                *string
	LogoPublicId           *string
	VideoUrl               *string
	VideoPublicId          *string
	ImageUrls              []string
	PreviousImagePublicIds []string
	NewImagePublicIds      []string
}
