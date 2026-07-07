package pkg

import (
	"errors"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

var userTypes []string = []string{"entrepreneur", "other", "employee", "student"}

func getRegxForSpaceAndCharacters() *regexp.Regexp {
	return regexp.MustCompile(`^[a-zA-Z\s]+$`)
}

func ValidateUserFullname(fullname string) error {

	if len(fullname) < 3 {
		return errors.New("Fullname should contain atleast 3 characters")
	}

	reg := getRegxForSpaceAndCharacters()

	correct := reg.MatchString(fullname)

	if !correct {
		return errors.New("Fullname should not contain any special character, symbols and digits")
	}

	return nil
}

func ValidateUserEmail(email string) bool {
	reg := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

	return reg.MatchString(email)
}

func ValidateCountryCode(code string) bool {
	if len(code) == 0 || len(code) == 1 {
		return false
	}

	contains := strings.Contains(code, "+")

	if !contains {
		return false
	}

	return true
}

func ValidateCountryOrState(str string) bool {

	if len(str) < 2 {
		return false
	}

	if len(str) > 150 {
		return false
	}

	reg := getRegxForSpaceAndCharacters()

	return reg.MatchString(str)
}

func ValidatePhoneNumber(phoneNumber int) bool {

	if phoneNumber > 9999999999 || phoneNumber < 0000000000 {
		return false
	}

	strPhoneNo := strconv.Itoa(phoneNumber)

	if len(strPhoneNo) != 10 {
		return false
	}

	return true

}

func ValidateUserType(userType string) bool {

	contains := slices.Contains(userTypes, userType)

	if !contains {
		return false
	}

	return true
}
