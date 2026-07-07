package repository

import "espectro/entity"

type UserRepository interface {
	Register(user entity.UserEntity) (string, error)
}
