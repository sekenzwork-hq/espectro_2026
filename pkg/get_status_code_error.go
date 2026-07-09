package pkg

import (
	"errors"
	customerrors "espectro/custom_errors"
	"net/http"
)

func GetStatusCodeForError(err error) int {

	switch {
	case errors.As(err, &customerrors.ValidationErr) || errors.As(err, &customerrors.CredentialErr):
		return http.StatusNotAcceptable
	case errors.As(err, &customerrors.ServerErr):
		return http.StatusInternalServerError
	case errors.As(err, &customerrors.PermissionErr) || errors.As(err, &customerrors.AuthErr):
		return http.StatusUnauthorized

	default:
		return http.StatusOK
	}
}
