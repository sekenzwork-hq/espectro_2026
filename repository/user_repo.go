package repository

import (
	"espectro/entity"

	"github.com/google/uuid"
)

type UserRepository interface {
	Register(user entity.UserEntity) (uuid.UUID, error)
}
