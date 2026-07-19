package repository

import "mime/multipart"

type MediaServiceRepo interface {
	UploadFiles(files []*multipart.FileHeader, folderId string, override bool) ([]string, error)
	UploadFile(file *multipart.FileHeader, folderId string, override bool) (string, error)
	DeleteFile(globalFolderId string, endpointFolder string) error
	DeleteMutipleFiles(globalFolderId string, endpointFolders []string) error
	RetrieveAssetPublicIds(folderId string) ([]string, error)
	DeleteAssetsWithPublicIds(publicIds []string) error
}
