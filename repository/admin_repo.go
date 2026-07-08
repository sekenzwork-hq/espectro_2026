package repository

type AdminRepo interface {
	Login(email string, password string) (bool, error)
}
