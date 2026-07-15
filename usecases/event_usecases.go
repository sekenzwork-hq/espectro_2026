package usecases

import (
	customerrors "espectro/custom_errors"
	"espectro/entity"
	"espectro/pkg"
	"espectro/repository"
)

type EventUsecases struct {
	repo repository.EventRepo
}

func NewEventUsecases(repo repository.EventRepo) EventUsecases {
	return EventUsecases{repo: repo}
}

func (e EventUsecases) CreateEvent(event entity.EventFromJsonEntity) (entity.EventEntity, error) {

	emptyEntity := entity.EventEntity{}

	nameErr := pkg.ValidateSpectrumOrEventName(event.Name)

	if nameErr != nil {
		return emptyEntity, &customerrors.ValidationError{OrgError: "Invalid event name"}
	}

}
