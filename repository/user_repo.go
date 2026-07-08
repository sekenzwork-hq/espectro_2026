package repository

import (
	"espectro/entity"

	"github.com/google/uuid"
)

type UserRepository interface {
	RegisterUser(user entity.UserEntity) (uuid.UUID, error)
}
