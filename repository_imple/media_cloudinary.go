package repositoryimple

import (
	"context"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/google/uuid"
)

type MediaCloudinaryRepo struct {
	cld *cloudinary.Cloudinary
}

func NewMediaCloudinaryRepo(cld *cloudinary.Cloudinary) MediaCloudinaryRepo {
	return MediaCloudinaryRepo{cld: cld}
}

func (m MediaCloudinaryRepo) UploadFiles(files []*multipart.FileHeader, folderId string, override bool) ([]string, error) {

	urls := []string{}

	for i := range files {
		file := files[i]
		url, err := m.UploadFile(file, folderId, override)
		if err != nil {
			return []string{}, err
		}
		urls = append(urls, url)
	}

	return urls, nil
}
func (m MediaCloudinaryRepo) UploadFile(file *multipart.FileHeader, folderId string, override bool) (string, error) {

	ctx := context.TODO()

	var publicId string
	if override {
		publicId = folderId
	} else {
		publicId = uuid.NewString()
	}

	res, err := m.cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder:    folderId,
		Overwrite: &override,
		PublicID:  publicId,
	})
	return res.SecureURL, err
}

func (m MediaCloudinaryRepo) DeleteFolderWithFiles(globalFolderId string, endpointFolder string) error {

	ctx := context.TODO()

	_, delAssetErr := m.cld.Admin.DeleteAssetsByPrefix(ctx, admin.DeleteAssetsByPrefixParams{
		Prefix: []string{globalFolderId + endpointFolder},
	})

	if delAssetErr != nil {
		return delAssetErr
	}

	_, folderDelErr := m.cld.Admin.DeleteFolder(ctx, admin.DeleteFolderParams{
		Folder: globalFolderId,
	})
	if folderDelErr != nil {
		return folderDelErr
	}
	return nil

}

func (m MediaCloudinaryRepo) DeleteMutiFoldersWithFiles(globalFolderId string, endpointFolders []string) error {

	for i := range endpointFolders {

		endpointFolder := endpointFolders[i]
		err := m.DeleteFolderWithFiles(globalFolderId, endpointFolder)
		if err != nil {
			return err
		}
	}

	return nil

}
