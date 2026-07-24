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
	galleryRepo      repository.GalleryRepo
	mediaRepo        repository.MediaServiceRepo
	venueRepo        repository.VenueRepo
	eventGalleryRepo repository.EventGalleryRepo
	eventRepo        repository.EventRepo
}

func NewGalleryUsecases(
	galleryRepo repository.GalleryRepo,
	mediaRepo repository.MediaServiceRepo,
	venueRepo repository.VenueRepo,
	eventGalleryRepo repository.EventGalleryRepo,
	eventRepo repository.EventRepo,
) GalleryUsecases {
	return GalleryUsecases{
		galleryRepo:      galleryRepo,
		mediaRepo:        mediaRepo,
		venueRepo:        venueRepo,
		eventGalleryRepo: eventGalleryRepo,
		eventRepo:        eventRepo,
	}
}

func (g GalleryUsecases) CreateGallery(name string, images []*multipart.FileHeader) (entity.GalleryEntity, error) {

	emptyGallery := entity.GalleryEntity{}

	validationErr := g.validateGalleryData(nil, &name, images)
	if validationErr != nil {
		return emptyGallery, validationErr
	}

	imageUrls := []string{}
	galleryId := uuid.NewString()

	folderId := "gallery/" + galleryId
	if len(images) != 0 {
		urls, _, _, err := g.uploadGalleryImages(folderId, images)
		if err != nil {
			return emptyGallery, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
		}
		imageUrls = urls
	}

	gallery, insertionErr := g.galleryRepo.CreateGallery(entity.GalleryEntity{
		Id:        galleryId,
		Name:      name,
		ImageUrls: imageUrls,
	})

	if insertionErr != nil {
		if len(imageUrls) != 0 {
			go g.mediaRepo.DeleteFile("gallery/", galleryId)
		}
		return emptyGallery, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}

	return gallery, nil

}
func (g GalleryUsecases) UpdateGallery(galleryId string, name *string, images []*multipart.FileHeader) (entity.GalleryEntity, error) {

	emptyGallery := entity.GalleryEntity{}

	validationErr := g.validateGalleryData(&galleryId, name, images)
	if validationErr != nil {
		return emptyGallery, validationErr
	}

	var imageUrls pq.StringArray
	folderId := "gallery/" + galleryId

	prevPublicIds := []string{}
	newPublicIds := []string{}
	if len(images) != 0 {
		urls, prevIds, newIds, err := g.uploadGalleryImages(folderId, images)
		if err != nil {
			return emptyGallery, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
		}
		imageUrls = pq.StringArray(urls)
		prevPublicIds = prevIds
		newPublicIds = newIds
	}

	newGallery, updationErr := g.galleryRepo.UpdateGallery(entity.GalleryUpdateEntity{
		GallerId:  galleryId,
		Name:      name,
		ImageUrls: &imageUrls,
	})

	if updationErr != nil {
		if len(imageUrls) != 0 {
			go g.mediaRepo.DeleteAssetsWithPublicIds(newPublicIds)
		}
		if errors.Is(updationErr, gorm.ErrRecordNotFound) {
			return emptyGallery, &customerrors.NotFoundError{OrgError: "Gallery does not exist"}
		}
		return emptyGallery, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}

	go g.mediaRepo.DeleteAssetsWithPublicIds(prevPublicIds)

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

func (g GalleryUsecases) RetrieveGalleries(limit int, page int) ([]entity.GalleryEntity, error) {

	offset := pkg.GetOffset(limit, page)

	galleries, err := g.galleryRepo.RetrieveGalleries(limit, offset)
	if err != nil {
		return galleries, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}
	return galleries, nil

}

func (e GalleryUsecases) AddGalleryToEvent(eventAndGallery entity.EventGalleryEntity) error {

	if !pkg.ValidateUUID(eventAndGallery.EventId) {
		return &customerrors.ValidationError{OrgError: "Invalid event id"}
	} else if !pkg.ValidateUUID(eventAndGallery.GalleryId) {
		return &customerrors.ValidationError{OrgError: "Invalid gallery id"}
	}

	eventExists, eventErr := e.eventRepo.EventExists(eventAndGallery.EventId)

	if eventErr != nil {
		return &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	} else if !eventExists {
		return &customerrors.NotFoundError{OrgError: "Event does not exist"}
	}

	galleryExists, galleryErr := e.galleryRepo.GalleryExists(eventAndGallery.GalleryId)

	if galleryErr != nil {
		return &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	} else if !galleryExists {
		return &customerrors.NotFoundError{OrgError: "Gallery does not exist"}
	}

	insertionErr := e.eventGalleryRepo.AddGalleryToEvent(eventAndGallery)

	if insertionErr != nil {
		return &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}

	return nil
}

func (g GalleryUsecases) validateGalleryData(galleryId *string, name *string, images []*multipart.FileHeader) error {

	if galleryId != nil && !pkg.ValidateUUID(*galleryId) {
		return &customerrors.ValidationError{OrgError: "Invalid gallery id"}
	}

	if name != nil {
		if err := pkg.ValidateName(*name); err != nil {
			return &customerrors.ValidationError{OrgError: err.Error()}
		}
	}

	if len(images) > 10 {
		return &customerrors.ValidationError{OrgError: "Maximum number of images is 10"}
	}

	for i := range images {
		image := images[i]
		if image != nil && !pkg.ValidateImageSize(*image) {
			return &customerrors.ValidationError{OrgError: "Size of each image should be less than or equal to 2 MB"}
		}
		valid, err := pkg.ValidateImage(image)
		if err != nil {
			return &customerrors.ServerError{OrgError: "Something went wrong while operating"}
		} else if !valid {
			return &customerrors.ValidationError{OrgError: "Invalid image format"}
		}
	}
	return nil
}

func (g GalleryUsecases) uploadGalleryImages(folderId string, images []*multipart.FileHeader) (urls []string, prevPublicIds []string, newPublicIds []string, err error) {

	prevPublicIds, retrivalErr := g.mediaRepo.RetrieveAssetPublicIds(folderId)
	if retrivalErr != nil {
		return []string{}, []string{}, []string{}, retrivalErr
	}

	urls, newPublicIds, err = g.mediaRepo.UploadFiles(images, folderId, false)
	if err != nil {
		return []string{}, []string{}, []string{}, err
	}
	return urls, prevPublicIds, newPublicIds, nil
}
