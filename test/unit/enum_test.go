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

func TestEventEnums(t *testing.T) {

	validStatus := []string{"pending", "on_going", "scheduled", "cancelled"}

	for i := range validStatus {
		status := validStatus[i]
		if !enums.EventStatus(status).IsValid() {
			t.Error("Event status validation failed (correct) : ", status)
			return
		}
	}

	inValidStatus := []string{"fjaljflaj", "!@#$$", "pendinggg"}

	for i := range inValidStatus {
		status := inValidStatus[i]
		if enums.EventStatus(status).IsValid() {
			t.Error("Event status validation failed (incorrect) : ", status)
			return
		}
	}

	validTypes := []string{"speech", "concert", "competition"}

	for i := range validTypes {
		v := validTypes[i]
		if !enums.EventType(v).IsValid() {
			t.Error("Event type validation failed (correct) : ", v)
			return
		}
	}
	invalidTypes := []string{"spechhh", "!@##$$$", "fafaf"}

	for i := range invalidTypes {
		v := invalidTypes[i]
		if enums.EventType(v).IsValid() {
			t.Error("Event type validation failed (incorrect) : ", v)
			return
		}
	}

	validModes := []string{"online", "stadium", "hall"}

	for i := range validModes {
		v := validModes[i]
		if !enums.EventMode(v).IsValid() {
			t.Error("Event mode validation failed (correct) : ", v)
			return
		}
	}
	invalidModes := []string{"onlineee", "!@#$", "hallooo"}

	for i := range invalidModes {
		v := invalidModes[i]
		if enums.EventMode(v).IsValid() {
			t.Error("Event mode validation failed (incorrect) : ", v)
			return
		}
	}

}

func TestSpectrumEnums(t *testing.T) {

	validStatuses := []string{"pending", "on_going", "scheduled"}

	for i := range validStatuses {
		v := validStatuses[i]
		if !enums.SpectrumStatus(v).IsValid() {
			t.Error("Spectrum status validation failed (correct) : ", v)
			return
		}

	}
	invalidStatuses := []string{"!@#$", "on__going", "sccccheduled"}

	for i := range invalidStatuses {
		v := invalidStatuses[i]
		if enums.SpectrumStatus(v).IsValid() {
			t.Error("Spectrum status validation failed (incorrect) : ", v)
			return
		}

	}
}

func TestSponsorEnums(t *testing.T) {

	validTypes := []string{"amount", "equipment", "gallery_arrangements", "electronics", "others", "furnitures"}

	for i := range validTypes {
		v := validTypes[i]

		if !enums.SponsoredType(v).IsValid() {
			t.Error("Sponsor type validation failed (correct) : ", v)
			return
		}
	}
	invalidTypes := []string{"amooount", "equipmentttt", "gallery___arrangements"}

	for i := range invalidTypes {
		v := invalidTypes[i]

		if enums.SponsoredType(v).IsValid() {
			t.Error("Sponsor type validation failed (incorrect) : ", v)
			return
		}
	}
}
