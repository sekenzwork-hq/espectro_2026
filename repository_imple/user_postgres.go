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

func (u UserPostgresRepo) RegisterUser(user entity.UserDBCreateEntity) (entity.UserEntity, error) {

	var createdUser entity.UserEntity

	err := u.db.Raw(
		`
		INSERT INTO users
		(
		username,
		password,
		fullname,
		email,
		country,
		country_code,
		state,
		city,
		phone,
		user_type)

		VALUES (?,?,?,?,?,?,?,?,?,?)

		RETURNING

		id,
		fullname,
		username,
		email,
		country,
		country_code,
		state,
		city,
		phone,
		user_type,
		created_at;
		`,
		user.Username,
		user.Password,
		user.Fullname,
		user.Email,
		user.Country,
		user.CountryCode,
		user.State,
		user.City,
		user.PhoneNumber,
		user.Usertype,
	).Scan(&createdUser).Error

	return createdUser, err
}

func (u UserPostgresRepo) UserExists(userId string) (bool, error) {

	var exists bool
	err := u.db.Raw(
		`
		SELECT EXISTS (SELECT 1 FROM users WHERE id=?)
		`,
		userId,
	).Scan(&exists).Error

	return exists, err
}

func (u UserPostgresRepo) UsernameExists(username string) (bool, error) {

	var exists bool
	err := u.db.Raw(
		`
		SELECT EXISTS (SELECT 1 FROM users WHERE username=?)
		`,
		username,
	).Scan(&exists).Error

	return exists, err
}

func (u UserPostgresRepo) RetrieveUserCredByUsername(username string) (entity.UserDBCredentialsEntity, error) {

	var cred entity.UserDBCredentialsEntity

	err := u.db.Table("users").
		Select("id,username,password").
		Where("username=?", username).
		First(&cred).Error

	return cred, err

}
