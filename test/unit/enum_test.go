package unit

import (
	"espectro/enums"
	"testing"
)

func TestAdminUpdationModeEnums(t *testing.T) {

	validEnums := []enums.AdminUpdateMode{enums.EmailAndFullname, enums.FullnameOnly, enums.EmailOnly}

	for i := range validEnums {
		enum := validEnums[i]

		if !enum.IsValid() {
			t.Errorf("Admin updation enum validation failed (correct). Enum : %v\n", enum)
		}
	}

	invalidEnums := []enums.AdminUpdateMode{"Unknown", "   ", "!@#$***", "update_it"}

	for i := range invalidEnums {
		enum := invalidEnums[i]

		if enum.IsValid() {
			t.Errorf("Admin updation enum validation failed (incorrct). Enum : %v\n", enum)
		}
	}
}

func TestAdminRoleEnums(t *testing.T) {

	validEnums := []enums.AdminRole{enums.Member, enums.Leader, enums.Volunteer}

	for i := range validEnums {
		enum := validEnums[i]

		if !enum.IsValid() {
			t.Errorf("Admin role enum validation failed (correct). Enum : %v\n", enum)
		}
	}

	invalidEnums := []enums.AdminRole{"Unknown", "   ", "!@#$***", "update_it"}

	for i := range invalidEnums {
		enum := invalidEnums[i]

		if enum.IsValid() {
			t.Errorf("Admin updation role validation failed (incorrct). Enum : %v\n", enum)
		}
	}

	strEnums := []string{"leader", "member", "volunteer"}

	for i := range strEnums {

		strEn := strEnums[i]

		if _, ok := enums.AdminRole(strEn).ParseRole(strEn); !ok {
			t.Errorf("Admin role enum validation failed. Enum : %v\n", strEn)
		}
	}

}
