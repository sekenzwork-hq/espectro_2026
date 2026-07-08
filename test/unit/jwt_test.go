package unit

import (
	"espectro/pkg"
	"testing"
)

func TestJWTForAdmins(t *testing.T) {

	ids := []string{"872fdb52-a7a6-49d4-b5d5-821da774043a", "550e8400-e29b-41d4-a716-446655440000", "c35296c4-c233-4bb7-af0b-bb827391e917"}
	tokens := []string{}

	for i := range ids {
		id := ids[i]

		token, err := pkg.GenerateJWTForAdmin(id)

		if err != nil || len(token) == 0 {
			t.Errorf("JWT Token generation for admin failed. ID : %v, Token : %v, Error : %v\n", id, token, err)
			return
		}

		tokens = append(tokens, token)
	}

	for i := range tokens {
		token := tokens[i]
		orgId := ids[i]

		id, err := pkg.ParseJWTFromAdmin(token)

		if len(id) == 0 || err != nil || id != orgId {

			t.Errorf("JWT Token validation for admin failed. ID : %v, OrgID : %v, token : %v, Error : %v", id, orgId, token, err)
			return
		}
	}
}
