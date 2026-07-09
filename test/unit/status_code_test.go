package unit

import (
	customerrors "espectro/custom_errors"
	"espectro/pkg"
	"testing"
)

func TestStatusCodeGeneration(t *testing.T) {

	supportedErros := []error{
		&customerrors.AuthenticationError{},
		&customerrors.CredentialsError{},
		&customerrors.NotFoundError{},
		&customerrors.PermissionError{},
		&customerrors.ServerError{},
	}

	corresCodes := []int{
		401,
		406,
		404,
		401,
		500,
	}

	for i := range supportedErros {
		err := supportedErros[i]
		correctCode := corresCodes[i]

		genCode := pkg.GetStatusCodeForError(err)

		if genCode != correctCode {

			t.Errorf("Status code generation failed. Code got : %v, Correct code : %v, Error : %v\n", genCode, correctCode, err)
			return
		}
	}
}
