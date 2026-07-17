package usecases

import (
	customerrors "espectro/custom_errors"
	"espectro/database"
	"espectro/entity"
	"espectro/enums"
	"espectro/pkg"
	"espectro/repository"
	"mime/multipart"

	"github.com/google/uuid"
)

type SponsorUsecases struct {
	sponsorRepo      repository.SponsorRepo
	eventSponsorRepo repository.EventSponsors
	eventRepo        repository.EventRepo
	transation       database.TransactionManager
	mediaRepo        repository.MediaServiceRepo
}

func NewSponsorUsecases(
	sponsorRepo repository.SponsorRepo,
	eventSponsorRepo repository.EventSponsors,
	eventRepo repository.EventRepo,
	transaction database.TransactionManager,
	mediaRepo repository.MediaServiceRepo,
) SponsorUsecases {
	return SponsorUsecases{
		sponsorRepo:      sponsorRepo,
		eventSponsorRepo: eventSponsorRepo,
		transation:       transaction,
		eventRepo:        eventRepo,
		mediaRepo:        mediaRepo,
	}
}

func (s SponsorUsecases) CreateSponsor(name string, amount *float32, profileOrOrg *multipart.FileHeader, sponsoredType enums.SponsoredType, eventIds []string) (entity.SponsorEntity, error) {

	emptySponsor := entity.SponsorEntity{}
	if err := pkg.ValidateName(name); err != nil {
		return emptySponsor, &customerrors.ValidationError{OrgError: err.Error()}
	}

	if amount != nil && *amount < 0 {
		return emptySponsor, &customerrors.ValidationError{OrgError: "Amount should be 0 or greater"}
	}

	if !sponsoredType.IsValid() {
		return emptySponsor, &customerrors.ValidationError{OrgError: "Invalid sponsor type"}
	}

	if len(eventIds) > 10 {
		return emptySponsor, &customerrors.SizeError{OrgError: "Maximum number of event ids is 10"}
	}

	for i := range eventIds {
		eventId := eventIds[i]

		if !pkg.ValidateUUID(eventId) {
			return emptySponsor, &customerrors.ValidationError{OrgError: "One of the event ids is invalid"}
		}
	}

	if len(eventIds) != 0 {
		exist, existErr := s.eventRepo.CheckMultipleEventsExist(eventIds)
		if existErr != nil {
			return emptySponsor, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
		}
		if !exist {
			return emptySponsor, &customerrors.NotFoundError{OrgError: "Some of the event ids do not exist"}
		}
	}

	var profileOrOrgUrl *string
	sponsorId := uuid.New().String()

	if profileOrOrg != nil {
		url, err := s.mediaRepo.UploadFile(profileOrOrg, "sponsor/"+sponsorId)
		if err != nil {
			return emptySponsor, &customerrors.ServerError{OrgError: "Something went wrong while operating"}
		}
		profileOrOrgUrl = &url
	}

	var sponsor entity.SponsorEntity
	trmErr := s.transation.Run(func() error {

		defaultAmount := float32(0)
		if amount != nil {
			defaultAmount = *amount
		}
		newSponsor, insertionErr := s.sponsorRepo.CreateSponsor(entity.SponsorEntity{
			Id:              sponsorId,
			Name:            name,
			Amount:          defaultAmount,
			ProfileOrOrgUrl: profileOrOrgUrl,
			Sponsored:       sponsoredType,
		})

		if insertionErr != nil {
			return &customerrors.ServerError{OrgError: "Something went wrong while operating"}
		}

		if len(eventIds) != 0 {
			addingErr := s.eventSponsorRepo.AddSponsor(newSponsor.Id, eventIds)
			if addingErr != nil {
				return &customerrors.ServerError{OrgError: "Something went wrong while operating"}
			}
		}
		sponsor = newSponsor
		return nil
	})

	if trmErr != nil {
		s.mediaRepo.DeleteFolderWithFiles("sponsor/", sponsorId)
		return emptySponsor, trmErr
	}

	return sponsor, nil

}
