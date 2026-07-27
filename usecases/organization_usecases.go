package usecases

import (
	customerrors "espectro/custom_errors"
	"espectro/entity"
	"espectro/enums"
	"espectro/pkg"
	"espectro/repository"
	"mime/multipart"
	"strings"

	"github.com/google/uuid"
)

type OrganizationUsecases struct {
	organizationRepo repository.OrganizationRepo
	mediaRepo        repository.MediaServiceRepo
}

func NewOrganizationUsecases(organizationRepo repository.OrganizationRepo, mediaRepo repository.MediaServiceRepo) OrganizationUsecases {
	return OrganizationUsecases{organizationRepo: organizationRepo}
}

func (o OrganizationUsecases) CreateOrganization(organization entity.OrganizationCreateEntity) (entity.OrganizationEntity, error) {

	empty := entity.OrganizationEntity{}
	err := o.validateOrganizationData(
		nil,
		&organization.Fullname,
		&organization.Email,
		organization.Logo,
		&organization.Status,
		nil,
		&organization.Fullname,
		&organization.PhoneNumber,
		organization.WebsiteUrl,
		&organization.Industry,
		&organization.HeadQuarters,
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
		Id:           organizationId,
		Fullname:     organization.Fullname,
		Email:        organization.Email,
		LogoUrl:      logoUrl,
		Status:       enums.PendingOrganization,
		ApprovedBy:   nil,
		PhoneNumber:  organization.PhoneNumber,
		WebsiteUrl:   organization.WebsiteUrl,
		Industry:     organization.Industry,
		HeadQuarters: organization.HeadQuarters,
	})

	if err != nil {
		return empty, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}

	return newOrg, nil
}

func (o OrganizationUsecases) validateOrganizationData(
	id *string,
	name *string,
	email *string,
	logo *multipart.FileHeader,
	status *enums.OrganizationStatus,
	approvedBy *string,
	fullname *string,
	phoneNumber *string,
	websiteUrl *string,
	industry *string,
	headquarters *string,
) error {

	if id != nil && !pkg.ValidateUUID(*id) {
		return &customerrors.ValidationError{OrgError: "Invalid organization id"}
	}

	if name != nil {
		err := pkg.ValidateName(*name)
		if err != nil {
			return &customerrors.ValidationError{OrgError: err.Error()}
		}
	}

	if email != nil && !pkg.ValidateEmail(*email) {
		return &customerrors.ValidationError{OrgError: "Invalid email"}
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

	if fullname != nil {
		err := pkg.ValidateFullname(*fullname)
		if err != nil {
			return &customerrors.ValidationError{OrgError: err.Error()}
		}
	}

	if phoneNumber != nil && !pkg.ValidatePhoneNumber(*phoneNumber) {
		return &customerrors.ValidationError{OrgError: "Invalid phone number"}
	}

	if websiteUrl != nil {
		err := pkg.ValidateUrl(*websiteUrl, "website url")
		return &customerrors.ValidationError{OrgError: err.Error()}
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

	return nil

}
