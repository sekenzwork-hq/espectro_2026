package pkg

import "golang.org/x/crypto/bcrypt"

func EncryptPassword(orgPass string) (string, error) {

	bytes, err := bcrypt.GenerateFromPassword([]byte(orgPass), 10)

	return string(bytes), err
}

func CompareHashedPass(orgPass string, hashedPass string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(orgPass))

	return err == nil
}
