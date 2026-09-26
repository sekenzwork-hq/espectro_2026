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
	userRepo         repository.UserRepository
}

func NewExhibitionUsecases(
	exhibitionRepo repository.ExhibitionRepo,
	mediaRepo repository.MediaServiceRepo,
	eventRepo repository.EventRepo,
	orgainazationRepo repository.OrganizationRepo,
	adminRepo repository.AdminRepo,
	staffRepo repository.StaffRepo,
	userRepo repository.UserRepository,

) ExhibitionUsecases {
	return ExhibitionUsecases{
		exhibitionRepo:   exhibitionRepo,
		mediaRepo:        mediaRepo,
		eventRepo:        eventRepo,
		organizationRepo: orgainazationRepo,
		adminRepo:        adminRepo,
		staffRepo:        staffRepo,
		userRepo:         userRepo,
	}
}

func (e ExhibitionUsecases) CreateExhibition(exhibition entity.ExhibitionRawEntity) (entity.ExhibitionDBRetrieveEntityFromUserSide, error) {

	empty := entity.ExhibitionDBRetrieveEntityFromUserSide{}

	validationError := e.validateExhibitionData(nil, &exhibition.EventId, exhibition.UserId, nil, &exhibition.Category, &exhibition.OrganizationId, nil, nil, nil, nil, &exhibition.ItemTitle, exhibition.ItemImages, &exhibition.ItemDescription, nil)

	if validationError != nil {
		return empty, validationError
	}

	eventExists, eventCheckingError := e.eventRepo.EventExists(exhibition.EventId)

	if eventCheckingError != nil {
		return empty, &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
	} else if !eventExists {
		return empty, &customerrors.NotFoundError{DisplayError: "Event does not exist"}
	}

	orgExists, orgCheckingError := e.organizationRepo.OrganizationExists(exhibition.OrganizationId)

	if orgCheckingError != nil {
		return empty, &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
	} else if !orgExists {
		return empty, &customerrors.NotFoundError{DisplayError: "Organization does not exist"}
	}

	exhibitionId := uuid.NewString()

	folderId := "exhibition/" + exhibitionId

	urls, publicIds, uploadingErr := e.mediaRepo.UploadFiles(
		exhibition.ItemImages,
		folderId,
		false,
	)

	if uploadingErr != nil {
		return empty, &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
	}

	status := enums.PendingExhibition
	pqUrls := pq.StringArray(urls)

	createdExhibition, creationError := e.exhibitionRepo.CreateExhibition(entity.ExhibitionDBInputEntity{
		Id:              &exhibitionId,
		EventId:         &exhibition.EventId,
		UserId:          exhibition.UserId,
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

func (e ExhibitionUsecases) UpdateExhibitionFromUserSide(exhibitionId string, newExhibition entity.ExhibitionRawUpdateEntity) (entity.ExhibitionDBRetrieveEntityFromUserSide, error) {

	empty := entity.ExhibitionDBRetrieveEntityFromUserSide{}
	validationError := e.validateExhibitionData(
		&exhibitionId,
		newExhibition.EventId,
		newExhibition.UserId,
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
			return empty, &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
		} else if !eventExists {
			return empty, &customerrors.NotFoundError{DisplayError: "Event does not exist"}
		}
	}

	if newExhibition.OrganizationId != nil {
		orgExists, orgCheckingError := e.organizationRepo.OrganizationExists(*newExhibition.OrganizationId)

		if orgCheckingError != nil {
			return empty, &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
		} else if !orgExists {
			return empty, &customerrors.NotFoundError{DisplayError: "Organization does not exist"}
		}
	}

	folderId := "exhibition/" + exhibitionId
	var previousPublicIds []string
	var newPublicIds []string
	var newImageUrls *pq.StringArray

	if len(newExhibition.ItemImages) != 0 {

		pubIds, retrievalError := e.mediaRepo.RetrieveAssetPublicIds(folderId, 10)

		if retrievalError != nil {
			return empty, &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
		}

		previousPublicIds = pubIds

		urls, newPubIds, uploadingErr := e.mediaRepo.UploadFiles(newExhibition.ItemImages, folderId, false)
		newPublicIds = newPubIds

		imageUrls := pq.StringArray(urls)
		newImageUrls = &imageUrls

		if uploadingErr != nil {
			return empty, &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
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
			return empty, &customerrors.NotFoundError{DisplayError: "Exhibitions does not exist"}
		}

		return empty, &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
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

func (e ExhibitionUsecases) UpdateExhibitionFromAdminSide(exhibitionId string, newExhibition entity.ExhibitionRawUpdateEntity) (entity.ExhibitionDBRetrieveEntityFromAdminSide, error) {

	empty := entity.ExhibitionDBRetrieveEntityFromAdminSide{}

	validationError := e.validateExhibitionData(
		&exhibitionId,
		nil,
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
			return empty, &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
		}

		if !exists {
			return empty, &customerrors.NotFoundError{DisplayError: "Approved admin does not exist"}
		}
	}

	if newExhibition.AssignedStaffId != nil {

		exists, err := e.staffRepo.CheckStaffExists(*newExhibition.AssignedStaffId)

		if err != nil {
			return empty, &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
		}

		if !exists {
			return empty, &customerrors.NotFoundError{DisplayError: "Assigned staff does not exist"}
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
			return empty, &customerrors.NotFoundError{DisplayError: "Exhibitions does not exist"}
		} else {
			return empty, &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
		}
	}

	return updatedExhibition, nil

}

func (e ExhibitionUsecases) DeleteExhibition(exhibitionId string) error {

	if !pkg.ValidateUUID(exhibitionId) {
		return &customerrors.ValidationError{DisplayError: "Invalid exhibition id"}
	}

	deletionError := e.exhibitionRepo.DeleteExhibition(exhibitionId)

	if deletionError != nil {

		fmt.Println("Deletion error : ", deletionError)
		if errors.Is(deletionError, gorm.ErrRecordNotFound) {
			return &customerrors.NotFoundError{DisplayError: "Exhibition does not exist"}
		} else {
			return &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
		}

	}

	go func() {
		e.mediaRepo.DeleteFile("exhibition/", exhibitionId)
	}()

	return nil
}

func (e ExhibitionUsecases) RetrieveExhibitionFromUserSide(userId string, page int) ([]entity.ExhibitionDBRetrieveEntityFromUserSide, error) {

	empty := []entity.ExhibitionDBRetrieveEntityFromUserSide{}

	isCorrect := pkg.ValidateUUID(userId)

	if !isCorrect {
		return empty, &customerrors.ValidationError{DisplayError: "Invalid user id"}
	}

	offset := pkg.GetOffset(50, page)

	exhibitions, err := e.exhibitionRepo.RetrieveExhibitionFromUserSide(userId, offset, page)

	if err != nil {
		return empty, &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
	}

	return exhibitions, err
}

func (e ExhibitionUsecases) RetrieveExhibitionFromAdminSide(page int) ([]entity.ExhibitionDBRetrieveEntityFromAdminSide, error) {

	offset := pkg.GetOffset(50, page)

	exhibitions, err := e.exhibitionRepo.RetrieveExhibitionFromAdminSide(offset, page)

	if err != nil {
		return []entity.ExhibitionDBRetrieveEntityFromAdminSide{}, &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
	}

	return exhibitions, nil
}

func (e ExhibitionUsecases) validateExhibitionData(
	id *string,
	eventId *string,
	userId *string,
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
		return &customerrors.ValidationError{DisplayError: "Invalid exhibition id"}
	}

	if eventId != nil && !pkg.ValidateUUID(*eventId) {
		return &customerrors.ValidationError{DisplayError: "Invalid event id"}
	}

	if userId != nil && !pkg.ValidateUUID(*userId) {
		return &customerrors.ValidationError{DisplayError: "Invalid user id"}
	}

	if category != nil {

		if len(*category) > 100 {
			return &customerrors.ValidationError{DisplayError: "Category length should be less than or equal to 100"}
		}

		trimmed := strings.TrimSpace(*category)

		if len(trimmed) < 2 {
			return &customerrors.ValidationError{DisplayError: "Category length should be greater than or equal to 2"}
		}

		correct := regex.MatchString(trimmed)

		if !correct {
			return &customerrors.ValidationError{DisplayError: "Category shouldn't contain any special characters"}
		}

	}

	if organizationId != nil && !pkg.ValidateUUID(*organizationId) {
		return &customerrors.ValidationError{DisplayError: "Invalid organization id"}
	}

	if boothNumber != nil {

		if *boothNumber < 0 {
			return &customerrors.ValidationError{DisplayError: "Booth number should be positive"}
		} else if *boothNumber > 25 {
			return &customerrors.ValidationError{DisplayError: "Booth number should be less than or equal 25"}
		}
	}

	if availableSqft != nil && *availableSqft <= 0 {
		return &customerrors.ValidationError{DisplayError: "Square feet should be greater than or equal to 1"}
	}

	if assignedStaffId != nil && !pkg.ValidateUUID(*assignedStaffId) {
		return &customerrors.ValidationError{DisplayError: "Invalid staff id"}
	}

	if approvedBy != nil && !pkg.ValidateUUID(*approvedBy) {
		return &customerrors.ValidationError{DisplayError: "Invalid approved admin id"}
	}

	if tokenNumber != nil {
		if *tokenNumber < 0 {
			return &customerrors.ValidationError{DisplayError: "Token number should be positive"}
		} else if *tokenNumber > 25 {
			return &customerrors.ValidationError{DisplayError: "Token number should be less than or equal to 25"}
		}
	}

	if itemTitle != nil && len(*itemTitle) > 100 {

		return &customerrors.ValidationError{DisplayError: "Title length should be less than or equal to 100"}

	} else if itemTitle != nil {

		trimmed := strings.TrimSpace(*itemTitle)
		if len(trimmed) < 5 {
			return &customerrors.ValidationError{DisplayError: "Title length should be greater than or equal to 5"}
		}

		correct := regex.MatchString(trimmed)

		if !correct {
			return &customerrors.ValidationError{DisplayError: "Title shouldn't contain any special characters"}
		}

	}

	if itemDescription != nil && len(*itemDescription) > 500 {

		return &customerrors.ValidationError{DisplayError: "Description should be less than or equal to 500"}

	} else if itemDescription != nil {

		trimmed := strings.TrimSpace(*itemDescription)
		if len(trimmed) < 5 {
			return &customerrors.ValidationError{DisplayError: "Description length should be greater than or equal to 5"}
		}

		correct := regex.MatchString(trimmed)

		if !correct {
			return &customerrors.ValidationError{DisplayError: "Description shouldn't contain any special characters"}
		}
	}

	if status != nil && !status.IsValid() {
		return &customerrors.ValidationError{DisplayError: "Invalid status"}
	}

	if len(itemImages) == 0 {
		return nil
	}

	if len(itemImages) > 10 {
		return &customerrors.ValidationError{DisplayError: "Maximum number of images is 10"}
	}

	for i := range itemImages {

		image := itemImages[i]
		if image == nil {
			continue
		}

		if !pkg.ValidateImageSize(*image) {
			return &customerrors.SizeError{DisplayError: "Image size should be less than or equal to 2MB"}
		}

		correct, err := pkg.ValidateImage(image)

		if err != nil {
			return &customerrors.ServerError{DisplayError: "Something went wrong while operating"}
		} else if !correct {
			return &customerrors.ValidationError{DisplayError: "Invalid image"}
		}

	}

	return nil

}
