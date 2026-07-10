package repository

import "mime/multipart"

type MediaServiceRepo interface {
	UploadImages(imageFiles []*multipart.FileHeader, folderId string) ([]string, error)
}
