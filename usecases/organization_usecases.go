package usecases

import (
	"errors"
	customerrors "espectro/custom_errors"
	"espectro/entity"
	"espectro/enums"
	"espectro/pkg"
	"espectro/repository"
	"mime/multipart"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrganizationUsecases struct {
	organizationRepo repository.OrganizationRepo
	mediaRepo        repository.MediaServiceRepo
}

func NewOrganizationUsecases(organizationRepo repository.OrganizationRepo, mediaRepo repository.MediaServiceRepo) OrganizationUsecases {
	return OrganizationUsecases{organizationRepo: organizationRepo, mediaRepo: mediaRepo}
}

func (o OrganizationUsecases) CreateOrganization(organization entity.OrganizationCreateEntity) (entity.OrganizationEntity, error) {

	empty := entity.OrganizationEntity{}

	status := enums.PendingOrganization
	err := o.validateOrganizationData(
		nil,
		&organization.Email,
		organization.Logo,
		&status,
		nil,
		&organization.Fullname,
		&organization.PhoneNumber,
		organization.WebsiteUrl,
		&organization.Industry,
		&organization.HeadQuarters,
		&organization.OrganizationName,
	)
	if err != nil {
		return empty, err
	}

	var logoUrl *string
	organizationId := uuid.NewString()
	if organization.Logo != nil {
		url, _, err := o.mediaRepo.UploadFile(organization.Logo, "organization/"+organizationId, true)
		if err != nil {
			return empty, &customerrors.ServerError{OrgError: "Something went wrong while operatinsg"}
		}

		logoUrl = &url
	}

	newOrg, err := o.organizationRepo.CreateOrganization(entity.OrganizationEntity{
		Id:               organizationId,
		Fullname:         organization.Fullname,
		Email:            organization.Email,
		LogoUrl:          logoUrl,
		Status:           enums.PendingOrganization,
		ApprovedBy:       nil,
		PhoneNumber:      organization.PhoneNumber,
		WebsiteUrl:       organization.WebsiteUrl,
		Industry:         organization.Industry,
		HeadQuarters:     organization.HeadQuarters,
		OrganizationName: organization.OrganizationName,
	})

	if err != nil {
		return empty, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}

	return newOrg, nil
}

func (o OrganizationUsecases) UpdateOrganization(id string,
	email *string,
	logo *multipart.FileHeader,
	fullname *string,
	phoneNumber *string,
	websiteUrl *string,
	industry *string,
	headquarters *string,
	organizationName *string,
) (entity.OrganizationEntity, error) {

	empty := entity.OrganizationEntity{}
	status := enums.PendingOrganization
	validationErr := o.validateOrganizationData(&id, email, logo, &status, nil, fullname, phoneNumber, websiteUrl, industry, headquarters, organizationName)
	if validationErr != nil {
		return empty, validationErr
	}

	var prevPublicId *string
	var newPublicId *string
	var logoUrl *string

	folderId := "organization/" + id
	if logo != nil {
		oldPubIds, err := o.mediaRepo.RetrieveAssetPublicIds(folderId)
		if err != nil {
			return empty, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
		}
		if len(oldPubIds) != 0 {
			prevPublicId = &oldPubIds[0]
		}

		url, newPubId, err := o.mediaRepo.UploadFile(logo, folderId, true)
		if err != nil {
			return empty, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
		}
		newPublicId = &newPubId
		logoUrl = &url

	}

	newOrganization, err := o.organizationRepo.UpdateOrganization(entity.OrganizationUpdateEntity{
		Id:               id,
		Email:            email,
		LogoUrl:          logoUrl,
		Status:           &status,
		Fullname:         fullname,
		PhoneNumber:      phoneNumber,
		WebsiteUrl:       websiteUrl,
		Industry:         industry,
		HeadQuarters:     headquarters,
		OrganizationName: organizationName,
	})

	go func() {
		if err != nil && newPublicId != nil {
			o.mediaRepo.DeleteAssetsWithPublicIds([]string{*newPublicId})
		} else if prevPublicId != nil {
			o.mediaRepo.DeleteAssetsWithPublicIds([]string{*prevPublicId})
		}
	}()

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return empty, &customerrors.NotFoundError{OrgError: "Organization does not exist"}
	} else if err != nil {
		return empty, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}

	return newOrganization, nil
}

func (o OrganizationUsecases) DeleteOrganization(id string) error {

	if !pkg.ValidateUUID(id) {
		return &customerrors.ValidationError{OrgError: "Invalid organization_id"}
	}

	err := o.organizationRepo.DeleteOrganization(id)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &customerrors.NotFoundError{OrgError: "Organization does not exist"}
	} else if err != nil {
		return &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}

	return nil
}

func (o OrganizationUsecases) RetrieveOrganizations(limit int, page int) ([]entity.OrganizationEntity, error) {

	offset := pkg.GetOffset(limit, page)

	organizations, err := o.organizationRepo.RetrieveOrganizations(limit, offset)

	if err != nil {
		return organizations, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}

	return organizations, nil
}

func (o OrganizationUsecases) validateOrganizationData(
	id *string,
	email *string,
	logo *multipart.FileHeader,
	status *enums.OrganizationStatus,
	approvedBy *string,
	fullname *string,
	phoneNumber *string,
	websiteUrl *string,
	industry *string,
	headquarters *string,
	organizationName *string,
) error {
	if id != nil && !pkg.ValidateUUID(*id) {
		return &customerrors.ValidationError{OrgError: "Invalid organization id"}
	}
	if email != nil && !pkg.ValidateEmail(*email) {
		return &customerrors.ValidationError{OrgError: "Invalid email"}
	}

	if fullname != nil {
		err := pkg.ValidateFullname(*fullname)
		if err != nil {
			return &customerrors.ValidationError{OrgError: err.Error()}
		}
	}

	if logo != nil {
		correct, err := pkg.ValidateImage(logo)

		if err != nil {
			return &customerrors.ServerError{OrgError: "Something went wrong while operatinsg"}
		} else if !correct {
			return &customerrors.ValidationError{OrgError: "Invalid image format"}
		}
	}
	if websiteUrl != nil {
		err := pkg.ValidateUrl(*websiteUrl, "website url")
		if err != nil {
			return &customerrors.ValidationError{OrgError: err.Error()}
		}
	}

	if status != nil && !status.IsValid() {
		return &customerrors.ValidationError{OrgError: "Invalid organization status"}
	}

	if approvedBy != nil && !pkg.ValidateUUID(*approvedBy) {
		return &customerrors.ValidationError{OrgError: "Invalid approved id"}
	}

	if phoneNumber != nil && !pkg.ValidatePhoneNumber(*phoneNumber) {
		return &customerrors.ValidationError{OrgError: "Invalid phone number"}
	}

	if industry != nil {
		if len(*industry) > 100 {
			return &customerrors.ValidationError{OrgError: "Industry length should be less than or equal to 100"}
		}

		str := strings.TrimSpace(*industry)
		if len(str) < 2 {
			return &customerrors.ValidationError{OrgError: "Industry length should be atleast 3"}
		}
	}
	if headquarters != nil {
		if len(*headquarters) > 200 {
			return &customerrors.ValidationError{OrgError: "Headquarters length should be less than or equal to 200"}
		}

		str := strings.TrimSpace(*headquarters)
		if len(str) < 10 {
			return &customerrors.ValidationError{OrgError: "Headquarters length should be atleast 10"}
		}
	}

	if organizationName != nil {
		if len(*organizationName) < 2 {
			return &customerrors.ValidationError{OrgError: "Organization name length should be atleast 2"}
		} else if len(*organizationName) > 150 {
			return &customerrors.ValidationError{OrgError: "Organization name length should be less than or equal to 150"}
		}
	}

	return nil

}
