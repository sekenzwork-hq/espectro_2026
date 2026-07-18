package repository

import "mime/multipart"

type MediaServiceRepo interface {
	UploadFiles(files []*multipart.FileHeader, folderId string, override bool) ([]string, error)
	UploadFile(file *multipart.FileHeader, folderId string, override bool) (string, error)
	DeleteFile(globalFolderId string, endpointFolder string) error
	DeleteMutipleFiles(globalFolderId string, endpointFolders []string) error
	MoveFiles(oldFolderId string, newFolderId string) error
}
