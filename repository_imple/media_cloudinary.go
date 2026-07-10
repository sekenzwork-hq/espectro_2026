package repositoryimple

import (
	"context"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
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

	ctx := context.Background()

	res, err := m.cld.Upload.Upload(ctx, imageFile, uploader.UploadParams{
		Folder: folderId,
	})

	return res.SecureURL, err
}

func (m MediaCloudinaryRepo) DeleteFolderWithFiles(folderId string) error {
	return nil
}
