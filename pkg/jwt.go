package pkg

import (
	customerrors "espectro/custom_errors"
	"fmt"
	"os"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWTForAdmin(id string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"admin_id": id,
	})

	strToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRETE")))

	return strToken, err

}

func GenerateJWTForUser(userId string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"user_id": userId,
	})

	strToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRETE")))

	return strToken, err
}

func ParseJWTFromAdmin(token string) (string, error) {

	invalidTokenErr := &customerrors.ValidationError{DisplayError: "Invalid token"}

	parsedToken, parseErr := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS512 {
			return nil, invalidTokenErr
		}
		return []byte(os.Getenv("JWT_SECRETE")), nil
	})

	if parseErr != nil {
		logger.Info("Parsing token error : ", parseErr)
		return "", parseErr
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if !ok {
		return "", invalidTokenErr
	}

	fmt.Println("Claims : ", claims)

	id, idOk := claims["admin_id"]

	if !idOk {
		return "", invalidTokenErr
	}

	idStr, isIdStr := id.(string)

	if !isIdStr {
		return "", invalidTokenErr
	} else if len(idStr) != 36 {
		return "", invalidTokenErr
	} else {
		return idStr, nil
	}

}

func ParseJWTFromUser(token string) (string, error) {

	invalidTokenErr := &customerrors.ValidationError{DisplayError: "Invalid token"}

	parsedToken, parsingErr := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS512 {
			return nil, invalidTokenErr
		}
		return []byte(os.Getenv("JWT_SECRETE")), nil
	})

	if parsingErr != nil {
		return "", parsingErr
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if !ok {
		return "", invalidTokenErr
	}

	id, idOk := claims["user_id"]

	if !idOk {
		return "", invalidTokenErr
	}

	idStr, isString := id.(string)

	if !isString {
		return "", invalidTokenErr
	} else if len(idStr) != 36 {
		return "", invalidTokenErr
	} else {
		return idStr, nil
	}

}
