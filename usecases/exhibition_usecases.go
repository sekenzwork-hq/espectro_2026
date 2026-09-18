package usecases

import (
	customerrors "espectro/custom_errors"
	"espectro/entity"
	"espectro/enums"
	"espectro/pkg"
	"espectro/repository"
	"fmt"
	"mime/multipart"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

type ExhibitionUsecases struct {
	exhibitionRepo   repository.ExhibitionRepo
	eventRepo        repository.EventRepo
	organizationRepo repository.OrganizationRepo
	mediaRepo        repository.MediaServiceRepo
}

func NewExhibitionUsecases(
	exhibitionRepo repository.ExhibitionRepo,
	mediaRepo repository.MediaServiceRepo,
	eventRepo repository.EventRepo,
	orgainazationRepo repository.OrganizationRepo,
) ExhibitionUsecases {
	return ExhibitionUsecases{
		exhibitionRepo:   exhibitionRepo,
		mediaRepo:        mediaRepo,
		eventRepo:        eventRepo,
		organizationRepo: orgainazationRepo,
	}
}

func (e ExhibitionUsecases) CreateExhibition(exhibition entity.ExhibitionCreateEntity) (entity.ExhibitionEntity, error) {

	empty := entity.ExhibitionEntity{}
	status := enums.PendingExhibition
	validationError := e.validateExhibitionData(nil, &exhibition.EventId, nil, &exhibition.Category, &exhibition.OrganizationId, nil, nil, nil, nil, &exhibition.ItemTitle, exhibition.ItemImages, &exhibition.ItemDescription, &status)

	if validationError != nil {
		return empty, validationError
	}

	eventExists, eventCheckingError := e.eventRepo.EventExists(exhibition.EventId)

	if eventCheckingError != nil {
		return empty, &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
	} else if !eventExists {
		return empty, &customerrors.NotFoundError{OrgError: "Event does not exist"}
	}

	orgExists, orgCheckingError := e.organizationRepo.OrganizationExists(exhibition.OrganizationId)

	if orgCheckingError != nil {
		return empty, &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
	} else if !orgExists {
		return empty, &customerrors.NotFoundError{OrgError: "Organization does not exist"}
	}

	exhibitionId := uuid.NewString()

	folderId := "exhibition/" + exhibitionId

	urls, publicIds, uploadingErr := e.mediaRepo.UploadFiles(
		exhibition.ItemImages,
		folderId,
		false,
	)

	if uploadingErr != nil {
		return empty, &customerrors.ServerError{DisplayError: "Something went wrong while operating", OrgError: uploadingErr.Error()}
	}

	createdExhibition, creationError := e.exhibitionRepo.CreateExhibition(entity.ExhibitionDBCreateEntity{
		Id:              exhibitionId,
		EventId:         exhibition.EventId,
		Category:        exhibition.Category,
		OrganizationId:  exhibition.OrganizationId,
		ItemTitle:       exhibition.ItemTitle,
		ItemImageUrls:   urls,
		ItemDescription: exhibition.ItemDescription,
		Status:          enums.PendingExhibition,
	})

	if creationError != nil {
		go func() {
			deletionErr := e.mediaRepo.DeleteAssetsWithPublicIds(publicIds)
			if deletionErr != nil {
			}
		}()
		fmt.Println("Exhibiton insertion error : ", creationError)
		return empty, &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
	}

	return createdExhibition, nil

}

func (e ExhibitionUsecases) validateExhibitionData(
	id *string,
	eventId *string,
	tokenNumber *int,
	category *string,
	organizationId *string,
	boothNumber *int,
	availableSqft *float64,
	assignedStaffId *string,
	approvedBy *string,
	itemTitle *string,
	itemImages []*multipart.FileHeader,
	itemDescription *string,
	status *enums.ExhibitionStatus,
) error {

	regex := regexp.MustCompile(`^[a-zA-Z0-9\s\-]+$`)

	if id != nil && !pkg.ValidateUUID(*id) {
		return &customerrors.ValidationError{OrgError: "Invalid exhibition id"}
	}

	if eventId != nil && !pkg.ValidateUUID(*eventId) {
		return &customerrors.ValidationError{OrgError: "Invalid event id"}
	}

	if category != nil {

		if len(*category) > 100 {
			return &customerrors.ValidationError{OrgError: "Category length should be less than or equal to 100"}
		}

		trimmed := strings.TrimSpace(*category)

		if len(trimmed) < 2 {
			return &customerrors.ValidationError{OrgError: "Category length should be greater than or equal to 2"}
		}

		correct := regex.MatchString(trimmed)

		if !correct {
			return &customerrors.ValidationError{OrgError: "Category shouldn't contain any special characters"}
		}

	}

	if organizationId != nil && !pkg.ValidateUUID(*organizationId) {
		return &customerrors.ValidationError{OrgError: "Invalid organization id"}
	}

	if boothNumber != nil && *boothNumber < 0 {
		return &customerrors.ValidationError{OrgError: "Booth number should be positive"}
	}

	if availableSqft != nil && *availableSqft <= 0 {
		return &customerrors.ValidationError{OrgError: "Square feet should be greater than or equal to 1"}
	}

	if assignedStaffId != nil && !pkg.ValidateUUID(*assignedStaffId) {
		return &customerrors.ValidationError{OrgError: "Invalid staff id"}
	}

	if approvedBy != nil && !pkg.ValidateUUID(*approvedBy) {
		return &customerrors.ValidationError{OrgError: "Invalid approved admin id"}
	}

	if tokenNumber != nil && *tokenNumber < 0 {
		return &customerrors.ValidationError{OrgError: "Token number should be positive"}
	}

	if itemTitle != nil && len(*itemTitle) > 100 {

		return &customerrors.ValidationError{OrgError: "Title length should be less than or equal to 100"}

	} else if itemTitle != nil {

		trimmed := strings.TrimSpace(*itemTitle)
		if len(trimmed) < 5 {
			return &customerrors.ValidationError{OrgError: "Title length should be greater than or equal to 5"}
		}

		correct := regex.MatchString(trimmed)

		if !correct {
			return &customerrors.ValidationError{OrgError: "Title shouldn't contain any special characters"}
		}

	}

	if itemDescription != nil && len(*itemDescription) > 500 {

		return &customerrors.ValidationError{OrgError: "Description should be less than or equal to 500"}

	} else if itemDescription != nil {

		trimmed := strings.TrimSpace(*itemDescription)
		if len(trimmed) < 5 {
			return &customerrors.ValidationError{OrgError: "Description length should be greater than or equal to 5"}
		}

		correct := regex.MatchString(trimmed)

		if !correct {
			return &customerrors.ValidationError{OrgError: "Description shouldn't contain any special characters"}
		}
	}

	if status != nil && !status.IsValid() {
		return &customerrors.ValidationError{OrgError: "Invalid status"}
	}

	if len(itemImages) == 0 {
		return nil
	}

	if len(itemImages) > 10 {
		return &customerrors.ValidationError{OrgError: "Maximum number of images is 10"}
	}

	for i := range itemImages {

		image := itemImages[i]
		if image == nil {
			continue
		}

		if !pkg.ValidateImageSize(*image) {
			return &customerrors.SizeError{OrgError: "Image size should be less than or equal to 2MB"}
		}

		correct, err := pkg.ValidateImage(image)

		if err != nil {
			return &customerrors.ServerError{OrgError: "Something went wrong while operating"}
		} else if !correct {
			return &customerrors.ValidationError{OrgError: "Invalid image"}
		}

	}

	return nil

}
