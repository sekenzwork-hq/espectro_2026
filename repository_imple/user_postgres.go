package repositoryimple

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

func (u UserPostgresRepo) RegisterUser(user entity.UserEntity) (entity.UserEntity, error) {

	var createdUser entity.UserEntity

	err := u.db.Raw(
		`
		INSERT INTO users
		(fullname,
		email,
		country,
		state,
		city,
		phone,
		user_type)

		VALUES (?,?,?,?,?,?,?)

		RETURNING

		id,
		fullname,
		email,
		country,
		state,
		city,
		phone,
		user_type,
		created_at;
		`,
		user.Fullname,
		user.Email,
		user.Country,
		user.State,
		user.City,
		user.PhoneNumber,
		user.Usertype,
	).Scan(&createdUser).Error

	return createdUser, err
}

func (u UserPostgresRepo) CheckUserExists(userId string) (bool, error) {

	var exists bool
	err := u.db.Raw(
		`
		SELECT EXISTS (SELECT 1 FROM users WHERE id=?)
		`,
		userId,
	).Scan(&exists).Error

	return exists, err
}
