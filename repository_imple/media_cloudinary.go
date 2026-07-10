package repositoryimple

import (
	"context"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type MediaCloudinaryRepo struct {
	cld *cloudinary.Cloudinary
}

func NewMediaCloudinaryRepo(cld *cloudinary.Cloudinary) MediaCloudinaryRepo {
	return MediaCloudinaryRepo{cld: cld}
}

func (m MediaCloudinaryRepo) UploadFiles(imageFiles []*multipart.FileHeader, folderId string) ([]string, error) {

	urls := []string{}
	ctx := context.Background()

	for i := range imageFiles {
		image := imageFiles[i]
		res, err := m.cld.Upload.Upload(ctx, image, uploader.UploadParams{
			Folder: folderId,
		})

		if err != nil {
			return []string{}, err
		}

		urls = append(urls, res.SecureURL)
	}

	return urls, nil
}
func (m MediaCloudinaryRepo) UploadFile(imageFile *multipart.FileHeader, folderId string) (string, error) {

	ctx := context.TODO()

	res, err := m.cld.Upload.Upload(ctx, imageFile, uploader.UploadParams{
		Folder: folderId,
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
