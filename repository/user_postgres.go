package repository

import (
	"espectro/entity"

	"gorm.io/gorm"
)

type UserPostgresRepo struct {
	db *gorm.DB
}

func NewUserPostgresRepo(db *gorm.DB) UserPostgresRepo {
	return UserPostgresRepo{
		db: db,
	}
}

func (u *UserPostgresRepo) RegisterUser(user entity.UserEntity) (string, error) {
	obj := u.db.Create(&user)
	return user.Id, obj.Error
}
