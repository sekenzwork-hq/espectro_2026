package pkg

import (
	"errors"
	"regexp"
	"slices"
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

	if len(fullname) > 100 {
		return errors.New("Length of fullname should be less than or equal to 100")
	}

	reg := getRegxForSpaceAndCharacters()

	correct := reg.MatchString(strings.TrimSpace(fullname))

	if !correct {
		return errors.New("Fullname should not contain any special character, symbols and digits")
	}

	return nil
}

func ValidateUserEmail(email string) bool {

	if len(email) > 255 {
		return false
	}
	reg := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

	return reg.MatchString(strings.TrimSpace(email))
}

func ValidateCountryCode(code string) bool {
	if len(code) == 0 || len(code) == 1 || len(code) > 200 {
		return false
	}

	count := 0

	for i := range code {

		var r rune = rune(code[i])

		if r == '+' {
			count++
		}

		if count == 2 {
			return false
		}
	}

	if count != 1 {
		return false
	}

	return true
}

func ValidateCountryOrState(str string) bool {

	if len(str) > 1000 {
		return false
	} else if len(str) <= 1 {
		return false
	}

	str = strings.TrimSpace(str)

	if len(str) > 150 || len(str) <= 1 {
		return false
	}

	reg := regexp.MustCompile(`^[A-Za-z\s()]+$`)

	match := reg.MatchString(str)

	return match
}

func ValidateCity(city string) bool {

	if len(city) < 5 {
		return false
	} else if len(city) > 200 {
		return false
	}

	city = strings.TrimSpace(city)

	reg := regexp.MustCompile(`^[A-Za-z\s()]+$`)

	match := reg.MatchString(city)

	return match

}

func ValidatePhoneNumber(phoneNumber string) bool {

	if len(phoneNumber) > 15 {
		return false
	}
	phoneNumber = strings.TrimSpace(phoneNumber)
	if len(phoneNumber) != 10 {
		return false
	}

	reg := regexp.MustCompile(`^[0-9]{10}$`)

	return reg.MatchString(phoneNumber)

}

func ValidateUserType(userType string) bool {

	userType = strings.TrimSpace(userType)

	contains := slices.Contains(userTypes, userType)

	if !contains {
		return false
	}

	return true
}
