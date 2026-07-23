package usecases

import (
	"errors"
	customerrors "espectro/custom_errors"
	"espectro/entity"
	"espectro/pkg"
	"espectro/repository"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type GalleryUsecases struct {
	galleryRepo repository.GalleryRepo
	mediaRepo   repository.MediaServiceRepo
	venueRepo   repository.VenueRepo
}

func NewGalleryUsecases(galleryRepo repository.GalleryRepo, mediaRepo repository.MediaServiceRepo, venueRepo repository.VenueRepo) GalleryUsecases {
	return GalleryUsecases{galleryRepo: galleryRepo, mediaRepo: mediaRepo, venueRepo: venueRepo}
}

func (g GalleryUsecases) CreateGallery(name string, venueId string, images []*multipart.FileHeader) (entity.GalleryEntity, error) {

	emptyGallery := entity.GalleryEntity{}

	validationErr := g.validateGalleryData(nil, &venueId, &name, images)

	if validationErr != nil {
		return emptyGallery, validationErr
	}

	imageUrls := []string{}
	galleryId := uuid.NewString()

	folderId := "gallery/" + galleryId
	if len(images) != 0 {
		urls, err := g.uploadGalleryImages(folderId, images)
		if err != nil {
			return emptyGallery, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
		}
		imageUrls = urls
	}

	gallery, insertionErr := g.galleryRepo.CreateGallery(entity.GalleryEntity{
		Id:        galleryId,
		Name:      name,
		ImageUrls: imageUrls,
		VenueId:   venueId,
	})

	if insertionErr != nil {
		if len(imageUrls) != 0 {
			go g.mediaRepo.DeleteFile("gallery/", galleryId)
		}
		return emptyGallery, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}

	return gallery, nil

}
func (g GalleryUsecases) UpdateGallery(galleryId string, name *string, venueId *string, images []*multipart.FileHeader) (entity.GalleryEntity, error) {

	emptyGallery := entity.GalleryEntity{}

	validationErr := g.validateGalleryData(&galleryId, venueId, name, images)
	if validationErr != nil {
		return emptyGallery, validationErr
	}

	var imageUrls pq.StringArray
	folderId := "gallery/" + galleryId

	if len(images) != 0 {
		urls, err := g.mediaRepo.UploadFiles(images, folderId, true)
		if err != nil {
			return emptyGallery, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
		}
		imageUrls = pq.StringArray(urls)
	}

	newGallery, updationErr := g.galleryRepo.UpdateGallery(entity.GalleryUpdateEntity{
		GallerId:  galleryId,
		Name:      name,
		ImageUrls: &imageUrls,
		VenueId:   venueId,
	})

	if updationErr != nil {
		if len(imageUrls) != 0 {
			go g.mediaRepo.DeleteFile("gallery/", galleryId)
		}
		if errors.Is(updationErr, gorm.ErrRecordNotFound) {
			return emptyGallery, &customerrors.NotFoundError{OrgError: "Gallery does not exist"}
		}
		return emptyGallery, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}

	return newGallery, updationErr

}

func (g GalleryUsecases) DeleteGallery(galleryId string) error {

	if !pkg.ValidateUUID(galleryId) {
		return &customerrors.ValidationError{OrgError: "Invalid gallery id"}
	}

	err := g.galleryRepo.DeleteGallery(galleryId)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &customerrors.NotFoundError{OrgError: "Gallery does not exist"}
	} else if err != nil {
		return &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}

	return nil
}

func (g GalleryUsecases) RetrieveGalleries(limit int, page int) ([]entity.GalleryWithVenueEntity, error) {

	offset := pkg.GetOffset(limit, page)

	galleries, err := g.galleryRepo.RetrieveGalleries(limit, offset)
	if err != nil {
		return galleries, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}
	return galleries, nil

}

func (g GalleryUsecases) validateGalleryData(galleryId *string, venueId *string, name *string, images []*multipart.FileHeader) error {

	if galleryId != nil && !pkg.ValidateUUID(*galleryId) {
		return &customerrors.ValidationError{OrgError: "Invalid gallery id"}
	}
	if venueId != nil && !pkg.ValidateUUID(*venueId) {
		return &customerrors.ValidationError{OrgError: "Invalid venue id"}
	}

	if name != nil {
		if err := pkg.ValidateName(*name); err != nil {
			return &customerrors.ValidationError{OrgError: err.Error()}
		}
	}

	if len(images) > 10 {
		return &customerrors.ValidationError{OrgError: "Maximum number of images is 10"}
	}
	if venueId != nil {

		venueExists, checkingErr := g.venueRepo.CheckVenueExists(*venueId)

		if checkingErr != nil || !venueExists {
			if checkingErr != nil {
				return &customerrors.ServerError{OrgError: "Something went wrong while operating"}
			} else {
				return &customerrors.NotFoundError{OrgError: "Venue does not exist"}
			}
		}

	}
	for i := range images {
		image := images[i]
		if image != nil && !pkg.ValidateImageSize(*image) {
			return &customerrors.ValidationError{OrgError: "Size of each image should be less than or equal to 2 MB"}
		}
	}
	return nil
}

func (g GalleryUsecases) uploadGalleryImages(folderId string, images []*multipart.FileHeader) ([]string, error) {

	prevPublicIds, retrivalErr := g.mediaRepo.RetrieveAssetPublicIds(folderId)

	if retrivalErr != nil {
		return []string{}, retrivalErr
	}

	urls, uploadErr := g.mediaRepo.UploadFiles(images, folderId, true)

	if uploadErr != nil {
		return []string{}, uploadErr
	}

	go g.mediaRepo.DeleteAssetsWithPublicIds(prevPublicIds)

	return urls, nil
}
