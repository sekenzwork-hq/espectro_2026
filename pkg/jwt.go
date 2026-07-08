package pkg

import (
	customerrors "espectro/custom_errors"
	"os"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWTForAdmin(id string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"id": id,
	})

	strToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRETE")))

	return strToken, err

}

func ParseJWTFromAdmin(token string) (string, error) {

	invalidTokenErr := &customerrors.ValidationError{OrgError: "Invalid token"}

	parsedToken, parseErr := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS512 {
			return nil, invalidTokenErr
		}
		return []byte(os.Getenv("JWT_SECRETE")), nil
	})

	if parseErr != nil {
		logger.Info("Parsing token error : ", parseErr)
		return "", &customerrors.ServerError{OrgError: "Something went wrong while operating"}
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if !ok {
		return "", invalidTokenErr
	}

	id, idOk := claims["id"]

	if !idOk {
		return "", invalidTokenErr
	}

	_, isIdStr := id.(string)

	if !isIdStr {
		return "", invalidTokenErr
	} else if len(id.(string)) != 36 {
		return "", invalidTokenErr
	} else {
		return id.(string), nil
	}

}
