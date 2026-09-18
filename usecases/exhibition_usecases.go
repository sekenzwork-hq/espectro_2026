package usecases

import (
	"errors"
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
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type ExhibitionUsecases struct {
	exhibitionRepo   repository.ExhibitionRepo
	eventRepo        repository.EventRepo
	organizationRepo repository.OrganizationRepo
	mediaRepo        repository.MediaServiceRepo
	adminRepo        repository.AdminRepo
	staffRepo        repository.StaffRepo
}

func NewExhibitionUsecases(
	exhibitionRepo repository.ExhibitionRepo,
	mediaRepo repository.MediaServiceRepo,
	eventRepo repository.EventRepo,
	orgainazationRepo repository.OrganizationRepo,
	adminRepo repository.AdminRepo,
	staffRepo repository.StaffRepo,

) ExhibitionUsecases {
	return ExhibitionUsecases{
		exhibitionRepo:   exhibitionRepo,
		mediaRepo:        mediaRepo,
		eventRepo:        eventRepo,
		organizationRepo: orgainazationRepo,
		adminRepo:        adminRepo,
		staffRepo:        staffRepo,
	}
}

func (e ExhibitionUsecases) CreateExhibition(exhibition entity.ExhibitionRawEntity) (entity.ExhibitionDBRetrieveEntity, error) {

	empty := entity.ExhibitionDBRetrieveEntity{}

	validationError := e.validateExhibitionData(nil, &exhibition.EventId, nil, &exhibition.Category, &exhibition.OrganizationId, nil, nil, nil, nil, &exhibition.ItemTitle, exhibition.ItemImages, &exhibition.ItemDescription, nil)

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

	status := enums.PendingExhibition
	pqUrls := pq.StringArray(urls)

	createdExhibition, creationError := e.exhibitionRepo.CreateExhibition(entity.ExhibitionDBInputEntity{
		Id:              &exhibitionId,
		EventId:         &exhibition.EventId,
		Category:        &exhibition.Category,
		OrganizationId:  &exhibition.OrganizationId,
		ItemTitle:       &exhibition.ItemTitle,
		ItemImageUrls:   &pqUrls,
		ItemDescription: &exhibition.ItemDescription,
		Status:          &status,
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

func (e ExhibitionUsecases) UpdateExhibitionFromUserSide(exhibitionId string, newExhibition entity.ExhibitionRawUpdateEntity) (entity.ExhibitionDBRetrieveEntity, error) {

	empty := entity.ExhibitionDBRetrieveEntity{}
	validationError := e.validateExhibitionData(
		&exhibitionId,
		newExhibition.EventId,
		nil,
		newExhibition.Category,
		newExhibition.OrganizationId,
		nil,
		nil,
		nil,
		nil,
		newExhibition.ItemTitle,
		newExhibition.ItemImages,
		newExhibition.ItemDescription,
		nil)

	if validationError != nil {
		return empty, validationError
	}

	if newExhibition.EventId != nil {
		eventExists, eventCheckingError := e.eventRepo.EventExists(*newExhibition.EventId)

		if eventCheckingError != nil {
			return empty, &customerrors.ServerError{DisplayError: "Something went wrong while operating", OrgError: eventCheckingError.Error()}
		} else if !eventExists {
			return empty, &customerrors.NotFoundError{DisplayError: "Event does not exist"}
		}
	}

	if newExhibition.OrganizationId != nil {
		orgExists, orgCheckingError := e.organizationRepo.OrganizationExists(*newExhibition.OrganizationId)

		if orgCheckingError != nil {
			return empty, &customerrors.ServerError{DisplayError: "Something went wrong while operating", OrgError: orgCheckingError.Error()}
		} else if !orgExists {
			return empty, &customerrors.NotFoundError{OrgError: "Organization does not exist"}
		}
	}

	folderId := "exhibition/" + exhibitionId
	var previousPublicIds []string
	var newPublicIds []string
	var newImageUrls *pq.StringArray

	if len(newExhibition.ItemImages) != 0 {

		pubIds, retrievalError := e.mediaRepo.RetrieveAssetPublicIds(folderId, 10)

		if retrievalError != nil {
			return empty, &customerrors.ServerError{OrgError: retrievalError.Error(), DisplayError: "Something went wrong while operating"}
		}

		previousPublicIds = pubIds

		urls, newPubIds, uploadingErr := e.mediaRepo.UploadFiles(newExhibition.ItemImages, folderId, false)
		newPublicIds = newPubIds

		imageUrls := pq.StringArray(urls)
		newImageUrls = &imageUrls

		if uploadingErr != nil {
			return empty, &customerrors.ServerError{OrgError: uploadingErr.Error(), DisplayError: "Something went wrong while operating"}
		}
	}

	updatedExhibition, updationError := e.exhibitionRepo.UpdateExhibitionFromUserSide(exhibitionId, entity.ExhibitionDBInputEntity{
		EventId:         newExhibition.EventId,
		Category:        newExhibition.Category,
		OrganizationId:  newExhibition.OrganizationId,
		ItemTitle:       newExhibition.ItemTitle,
		ItemImageUrls:   (*pq.StringArray)(newImageUrls),
		ItemDescription: newExhibition.ItemDescription,
	})

	if updationError != nil {
		go func() {
			deletionErr := e.mediaRepo.DeleteAssetsWithPublicIds(newPublicIds)
			if deletionErr != nil {
				fmt.Println("Exhibition images deletion error : ", deletionErr)
			}
		}()

		if errors.Is(updationError, gorm.ErrRecordNotFound) {
			return empty, &customerrors.NotFoundError{OrgError: updationError.Error(), DisplayError: "Exhibition does not exist"}
		}

		return empty, &customerrors.ServerError{OrgError: updationError.Error(), DisplayError: "Something went wrong while operating"}
	} else {
		go func() {
			deletionErr := e.mediaRepo.DeleteAssetsWithPublicIds(previousPublicIds)
			if deletionErr != nil {
				fmt.Println("Exhibition images deletion error : ", deletionErr)
			}
		}()

		return updatedExhibition, nil
	}

}

func (e ExhibitionUsecases) UpdateExhibitionFromAdminSide(exhibitionId string, newExhibition entity.ExhibitionRawUpdateEntity) (entity.ExhibitionDBRetrieveEntity, error) {

	empty := entity.ExhibitionDBRetrieveEntity{}

	validationError := e.validateExhibitionData(
		&exhibitionId,
		nil,
		newExhibition.TokenNumber,
		nil,
		nil,
		newExhibition.BoothNumber,
		newExhibition.AvailableSqft,
		newExhibition.AssignedStaffId,
		newExhibition.ApprovedBy,
		nil,
		nil,
		nil,
		newExhibition.Status,
	)

	if validationError != nil {
		return empty, validationError
	}

	if newExhibition.ApprovedBy != nil {

		exists, err := e.adminRepo.CheckAdminExists(*newExhibition.ApprovedBy)

		if err != nil {
			return empty, &customerrors.ServerError{OrgError: err.Error(), DisplayError: "Something went wrong while operating"}
		}

		if !exists {
			return empty, &customerrors.NotFoundError{OrgError: "Exists returned false", DisplayError: "Approved admin does not exist"}
		}
	}

	if newExhibition.AssignedStaffId != nil {

		exists, err := e.staffRepo.CheckStaffExists(*newExhibition.AssignedStaffId)

		if err != nil {
			return empty, &customerrors.ServerError{OrgError: err.Error(), DisplayError: "Something went wrong while operating"}
		}

		if !exists {
			return empty, &customerrors.NotFoundError{OrgError: "Exists returned false", DisplayError: "Assigned staff does not exist"}
		}
	}

	updatedExhibition, updationError := e.exhibitionRepo.UpdateExhibitionFromAdminSide(exhibitionId, entity.ExhibitionDBInputEntity{
		TokenNumber:     newExhibition.TokenNumber,
		BoothNumber:     newExhibition.BoothNumber,
		AvailableSqft:   newExhibition.AvailableSqft,
		AssignedStaffId: newExhibition.AssignedStaffId,
		ApprovedBy:      newExhibition.ApprovedBy,
		Status:          newExhibition.Status,
	})

	if updationError != nil {

		if errors.Is(updationError, gorm.ErrRecordNotFound) {
			return empty, &customerrors.NotFoundError{DisplayError: "Exhibition does not exist", OrgError: updationError.Error()}
		} else {
			return empty, &customerrors.ServerError{DisplayError: "Something went wrong", OrgError: updationError.Error()}
		}
	}

	return updatedExhibition, nil

}

func (e ExhibitionUsecases) DeleteExhibition(exhibitionId string) error {

	if !pkg.ValidateUUID(exhibitionId) {
		return &customerrors.ValidationError{OrgError: "Invalid exhibition id"}
	}

	deletionError := e.exhibitionRepo.DeleteExhibition(exhibitionId)

	if deletionError != nil {

		fmt.Println("Deletion error : ", deletionError)
		if errors.Is(deletionError, gorm.ErrRecordNotFound) {
			return &customerrors.NotFoundError{OrgError: "Exhibition does not exist", DisplayError: "Exhibition does not exist"}
		} else {
			return &customerrors.ServerError{OrgError: deletionError.Error(), DisplayError: "Something went wrong while operating"}
		}

	}

	go func() {
		e.mediaRepo.DeleteFile("exhibition/", exhibitionId)
	}()

	return nil
}

func (e ExhibitionUsecases) validateExhibitionData(
	id *string,
	eventId *string,
	tokenNumber *int,
	category *string,
	organizationId *string,
	boothNumber *int,
	availableSqft *float32,
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

	if boothNumber != nil {

		if *boothNumber < 0 {
			return &customerrors.ValidationError{OrgError: "Booth number should be positive"}
		} else if *boothNumber > 25 {
			return &customerrors.ValidationError{OrgError: "Booth number should be less than or equal 25"}
		}
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

	if tokenNumber != nil {
		if *tokenNumber < 0 {
			return &customerrors.ValidationError{OrgError: "Token number should be positive"}
		} else if *tokenNumber > 25 {
			return &customerrors.ValidationError{OrgError: "Token number should be less than or equal to 25"}
		}
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
