package unit

import (
	"espectro/pkg"
	"testing"
)

func TestValidateFullname(t *testing.T) {

	correctSampleNames := []string{
		"Adhun U",
		"Adhun",
		"John",
		"Joe",
		"Kanaaran",
	}

	for i := range correctSampleNames {

		name := correctSampleNames[i]

		err := pkg.ValidateFullname(name)

		if err != nil {
			t.Error("Validation of fullname failed the test case (correct names) : ", name)
			return
		}

	}

	incorrectSampleNames := []string{"Ad", "AD", "Jo", "123", "@#$%", "Raju123", "Soman@#$", "", "   "}

	for i := range incorrectSampleNames {

		name := incorrectSampleNames[i]

		err := pkg.ValidateFullname(name)

		if err == nil {
			t.Error("Validation of fullname failed the test case (incorrect name) : ", name)
			return
		}
	}
}

func TestValidateEmail(t *testing.T) {

	correctEmails := []string{"adhun@gmail.com", "adhun1@gmail.com", "joe@gmail.com", "rajappan@yahoo.com", "soman@outlook.com", "someone@email.com"}

	for i := range correctEmails {

		email := correctEmails[i]

		correct := pkg.ValidateEmail(email)

		if !correct {
			t.Error("Validation of email failed the test case (correct) : ", email)
			return
		}
	}

	incorrectEmails := []string{
		"user#example.com",
		"user!@example.com",
		"user@example!.com",
		"user<>@example.com",
		"user@example,com",
		"user@example/com",
	}

	for i := range incorrectEmails {

		email := incorrectEmails[i]

		correct := pkg.ValidateEmail(email)

		if correct {
			t.Error("Validation of email failed the test case (incorrect) : ", email)
			return
		}
	}
}

func TestValidateCountryCode(t *testing.T) {

	correctCountryCodes := []string{"+91", "+1", "+44", "+61", "+39", "+351", "+966", "+1-268"}

	for i := range correctCountryCodes {

		code := correctCountryCodes[i]

		correct := pkg.ValidateCountryCode(code)

		if !correct {
			t.Error("Validation of country code failed the test case (correct) : ", code)
			return
		}
	}

	incorrectCountryCodes := []string{"-1", "@91", "##5", "  ", "++91", "10"}

	for i := range incorrectCountryCodes {

		code := incorrectCountryCodes[i]

		correct := pkg.ValidateCountryCode(code)

		if correct {
			t.Error("Validation of country code failed the test case (incorrect) : ", code)
			return
		}
	}
}

func TestValidCountryOrState(t *testing.T) {

	correctCountries := []string{"India", "US", "China", "Japan", "United Kingdom of Great Britain and Northern Ireland"}
	correctStates := []string{"Kerala", "Tokyo", "Washington", "Dadra and Nagar Haveli and Daman and Diu (Union Territory)"}

	for i := range correctCountries {

		country := correctCountries[i]

		correct := pkg.ValidateCountryOrState(country)

		if !correct {
			t.Error("Validation of country failed the test case (correct) : ", country)
			return
		}
	}

	for i := range correctStates {
		state := correctStates[i]

		correct := pkg.ValidateCountryOrState(state)

		if !correct {
			t.Error("Validation of state failed the test case (correct) : ", state)
			return
		}
	}

	incorrectCountriesOrStates := []string{"!@@##$%#$@", " ", "I", "J#$$%", "        ", "+-*/"}

	for i := range incorrectCountriesOrStates {
		countryOrState := incorrectCountriesOrStates[i]

		correct := pkg.ValidateCountryOrState(countryOrState)

		if correct {
			t.Error("Validation of country or state failed the test case (incorrect) : ", countryOrState)
			return
		}
	}
}

func TestValidatePhoneNumber(t *testing.T) {

	correctPhoneNumbers := []string{"9090909090", "8828758901", " 9090909090", " 9090909090 "}

	for i := range correctPhoneNumbers {

		number := correctPhoneNumbers[i]

		correct := pkg.ValidatePhoneNumber(number)

		if !correct {
			t.Error("Validation of phone number failed the test case (correct) : ", number)
			return
		}
	}

	incorrectPhoneNumbers := []string{"               ", "!@##909090", "12345678", "123", "!@#$^&**$#"}

	for i := range incorrectPhoneNumbers {
		number := incorrectPhoneNumbers[i]

		correct := pkg.ValidatePhoneNumber(number)

		if correct {
			t.Error("Validation of phone failed the test case (incorrect) : ", number)
			return
		}
	}
}

func TestValidateUserType(t *testing.T) {

	correctTypes := []string{"entrepreneur", "student", "employee", "other"}

	for i := range correctTypes {

		uType := correctTypes[i]

		correct := pkg.ValidateUserType(uType)

		if !correct {
			t.Error("Validation of user type failed the test case (correct) : ", uType)
			return
		}
	}

	incorrectTypes := []string{"!@@$#%#", "unknown", "        ", "//////", "||||"}

	for i := range incorrectTypes {

		uType := incorrectTypes[i]

		correct := pkg.ValidateUserType(uType)

		if correct {
			t.Error("Validation of user type failed the test case (incorrect) : ", uType)
			return
		}
	}
}

func TestValidateCity(t *testing.T) {

	correctCities := []string{
		"Thiruvananthapuram",
		"Chhatrapati Sambhajinagar",
		"Tiruchirappalli",
	}

	for i := range correctCities {

		city := correctCities[i]

		correct := pkg.ValidateCity(city)

		if !correct {
			t.Error("Validation of city failed the test case (correct) : ", city)
			return
		}
	}

	incorrectCities := []string{"1234", "cfff", "1234567777", "@!!@#$%^&&&&&&$"}

	for i := range incorrectCities {
		city := incorrectCities[i]

		correct := pkg.ValidateCity(city)

		if correct {
			t.Error("Validation of city failed the test case (incorrect) : ", city)
			return
		}
	}
}

func TestValidatePassword(t *testing.T) {

	correctPasswords := []string{"someone@#$123", "another@#456"}

	for i := range correctPasswords {

		pass := correctPasswords[i]

		err := pkg.ValidatePassword(pass)

		if err != nil {
			t.Errorf("Password validation failed (correct) : Error : %v, Password : %v", err, pass)
			return
		}
	}

	incorrectPass := []string{"12345678", "12345", "!@#$%", "someone@"}

	for i := range incorrectPass {

		pass := incorrectPass[i]

		err := pkg.ValidatePassword(pass)

		if err == nil {
			t.Errorf("Password validation failed (incorrect) : Error : %v, Password : %v", err, pass)
			return
		}
	}
}

func TestValidateAdminRoles(t *testing.T) {

	correctRoles := []string{"volunteer", "leader", "member"}

	for i := range correctRoles {
		role := correctRoles[i]

		correct := pkg.ValidateAdminRole(role)

		if !correct {
			t.Error("Admin role validation failed (correct) : ", role)
			return
		}
	}

	incorrectRoles := []string{"something", "1@###$$", "12344555"}
	for i := range incorrectRoles {
		role := incorrectRoles[i]

		correct := pkg.ValidateAdminRole(role)

		if correct {
			t.Error("Admin role validation failed (incorrect) : ", role)
			return
		}
	}
}
