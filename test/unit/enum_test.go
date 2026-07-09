package unit

import (
	"espectro/enums"
	"fmt"
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

	var e enums.AdminRole

	fmt.Println("Enum : ", e.IsValid())

}
