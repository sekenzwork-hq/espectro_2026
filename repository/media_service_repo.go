package repository

import "mime/multipart"

type MediaServiceRepo interface {
	UploadFiles(imageFiles []*multipart.FileHeader, folderId string) ([]string, error)
	UploadFile(imageFile *multipart.FileHeader, folderId string) (string, error)
	DeleteFolderWithFiles(folderId string) error
}
