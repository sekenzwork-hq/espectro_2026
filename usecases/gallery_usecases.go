package usecases

import (
	customerrors "espectro/custom_errors"
	"espectro/entity"
	"espectro/pkg"
	"espectro/repository"
	"mime/multipart"

	"github.com/google/uuid"
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
	if !pkg.ValidateUUID(venueId) {
		return emptyGallery, &customerrors.ValidationError{OrgError: "Invalid venue id"}
	}

	if err := pkg.ValidateName(name); err != nil {
		return emptyGallery, &customerrors.ValidationError{OrgError: err.Error()}
	}

	if len(images) > 10 {
		return emptyGallery, &customerrors.ValidationError{OrgError: "Maximum number of images is 10"}
	}

	venueExists, checkingErr := g.venueRepo.CheckVenueExists(venueId)

	if checkingErr != nil || !venueExists {
		if checkingErr != nil {
			return emptyGallery, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
		} else {
			return emptyGallery, &customerrors.NotFoundError{OrgError: "Venue does not exist"}
		}
	}

	for i := range images {
		image := images[i]
		if image != nil && !pkg.ValidateImageSize(*image) {
			return emptyGallery, &customerrors.ValidationError{OrgError: "Size of each image should be less than or equal to 2 MB"}
		}
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
		go g.mediaRepo.DeleteFile("gallery/", galleryId)
		return emptyGallery, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}

	return gallery, nil

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
