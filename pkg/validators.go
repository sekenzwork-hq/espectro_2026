package pkg

import (
	"errors"
	customerrors "espectro/custom_errors"
	"regexp"
	"slices"
	"strings"
)

func getRegxForSpaceAndCharacters() *regexp.Regexp {
	return regexp.MustCompile(`^[a-zA-Z\s]+$`)
}

func ValidateFullname(fullname string) error {

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

func ValidateEmail(email string) bool {

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
	} else if len(city) > 1000 {
		return false
	}

	city = strings.TrimSpace(city)

	if len(city) > 150 && len(city) < 4 {
		return false
	}

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

	userTypes := []string{"entrepreneur", "other", "employee", "student"}

	userType = strings.TrimSpace(userType)

	contains := slices.Contains(userTypes, userType)

	return contains
}

func ValidatePassword(password string) error {

	if len(password) < 6 {
		return &customerrors.ValidationError{OrgError: "Password length should be atleast 6"}
	} else if len(password) > 200 {
		return &customerrors.ValidationError{OrgError: "Password length should be less than or equal to 200"}
	}

	password = strings.TrimSpace(password)

	upperOrLowerRegex := regexp.MustCompile(`(?=.*[A-Za-z])`)
	numberRegex := regexp.MustCompile(`(?=.*\d)`)
	specialCharRegex := regexp.MustCompile(`(?=.*[^A-Za-z0-9])`)

	containsUpperOrLower := upperOrLowerRegex.MatchString(password)

	if !containsUpperOrLower {
		return &customerrors.ValidationError{OrgError: "Password should contain atleast one upper or lower case"}
	}

	containsNumber := numberRegex.MatchString(password)

	if !containsNumber {
		return &customerrors.ValidationError{OrgError: "Password should contain alteast one digit"}
	}

	containsSpecialChar := specialCharRegex.MatchString(password)

	if !containsSpecialChar {
		return &customerrors.ValidationError{OrgError: "Password should contain alteast one special character"}
	}

	return nil
}

func ValidateAdminRole(adminRole string) bool {

	roles := []string{"leader", "member", "volunteer"}

	contains := slices.Contains(roles, adminRole)

	return contains
}
