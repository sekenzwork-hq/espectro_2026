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

func (m MediaCloudinaryRepo) UploadFiles(files []*multipart.FileHeader, folderId string, override bool) (urls []string, publicIds []string, err error) {

	urls = []string{}
	publicIds = []string{}

	for i := range files {
		file := files[i]
		if file == nil {
			continue
		}
		url, publicId, err := m.UploadFile(file, folderId, override)
		if err != nil {
			return []string{}, []string{}, err
		}
		urls = append(urls, url)
		publicIds = append(publicIds, publicId)
	}

	return urls, publicIds, nil
}
func (m MediaCloudinaryRepo) UploadFile(file *multipart.FileHeader, folderId string, override bool) (url string, publicId string, err error) {

	ctx := context.TODO()

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
	return res.SecureURL, folderId + "/" + publicId, err
}

func (m MediaCloudinaryRepo) DeleteFile(globalFolderId string, endpointFolder string) error {

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

func (m MediaCloudinaryRepo) DeleteMutipleFiles(globalFolderId string, endpointFolders []string) error {

	for i := range endpointFolders {
		endpointFolder := endpointFolders[i]
		err := m.DeleteFile(globalFolderId, endpointFolder)
		if err != nil {
			return err
		}
	}

	return nil

}

func (m MediaCloudinaryRepo) RetrieveAssetPublicIds(folderId string, limit int) ([]string, error) {

	assets, err := m.cld.Admin.AssetsByAssetFolder(context.TODO(), admin.AssetsByAssetFolderParams{
		AssetFolder: folderId,
		MaxResults:  limit,
	})

	if err != nil {
		return []string{}, err
	}

	ids := []string{}
	for i := range assets.Assets {
		asset := assets.Assets[i]
		ids = append(ids, asset.PublicID)
	}

	return ids, nil
}

func (m MediaCloudinaryRepo) DeleteAssetsWithPublicIds(publicIds []string) error {

	_, err := m.cld.Admin.DeleteAssets(context.TODO(), admin.DeleteAssetsParams{
		PublicIDs: publicIds,
	})
	if err != nil {
		return err
	}
	return nil
}
