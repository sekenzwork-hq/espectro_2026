package repository

import (
	"espectro/entity"
)

type UserRepository interface {
	RegisterUser(user entity.UserEntity) (entity.UserEntity, error)
	CheckUserExists(userId string) (bool, error)
}
