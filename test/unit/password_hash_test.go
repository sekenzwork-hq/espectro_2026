package unit

import (
	"espectro/pkg"
	"testing"
)

func TestPasswordHash(t *testing.T) {

	newPasswords := []string{"someone@180", "u_1234$#", "joe$%#1234", "jfajf!@#$"}
	hashes := []string{}

	for i := range newPasswords {

		pass := newPasswords[i]

		hash, err := pkg.EncryptPassword(pass)

		if err != nil || len(hash) == 0 {
			t.Error("Password encryption failed : ", pass)
			return
		}

		hashes = append(hashes, hash)
	}

	for i := range hashes {

		hash := hashes[i]
		orgPass := newPasswords[i]

		correct := pkg.CompareHashedPass(orgPass, hash)

		if !correct {
			t.Errorf("Password comparing failed : org : %v, hash : %v", orgPass, hash)
			return
		}
	}
}
