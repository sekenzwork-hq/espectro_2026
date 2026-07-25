package repositoryimple

import (
	"espectro/entity"

	"github.com/google/uuid"
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

func (u UserPostgresRepo) RegisterUser(user entity.UserEntity) (uuid.UUID, error) {
	obj := u.db.
		Table("users").
		Create(&user)

	return user.Id, obj.Error
}

func (u UserPostgresRepo) UserExists(userId string) (bool, error) {

	var exists bool

	err := u.db.Raw(
		`
		SELECT EXISTS (SELECT 1 FROM users WHERE id=? AND deleted_at IS NULL)
		`,
		userId,
	).Scan(&exists).Error

	return exists, err
}
