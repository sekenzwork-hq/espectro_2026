package repository

import "mime/multipart"

type MediaServiceRepo interface {
	UploadFiles(files []*multipart.FileHeader, folderId string, override bool) ([]string, error)
	UploadFile(file *multipart.FileHeader, folderId string, override bool) (string, error)
	DeleteFolderWithFiles(globalFolderId string, endpointFolder string) error
	DeleteMutiFoldersWithFiles(globalFolderId string, endpointFolders []string) error
}
