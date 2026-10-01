package repository

import (
	"espectro/entity"
)

type UserRepository interface {
	RegisterUser(user entity.UserDBCreateEntity) (entity.UserEntity, error)
	UserExists(userId string) (bool, error)
	UsernameExists(username string) (bool, error)
	RetrieveUserCredByUsername(username string) (entity.UserDBCredentialsEntity, error)
}
