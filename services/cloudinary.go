package services

import (
	"github.com/cloudinary/cloudinary-go/v2"
)

func NewCloudinary(url string) (*cloudinary.Cloudinary, error) {

	cld, cldErr := cloudinary.NewFromURL(url)

	return cld, cldErr
}
