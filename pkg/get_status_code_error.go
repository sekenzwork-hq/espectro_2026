package pkg

import (
	"errors"
	customerrors "espectro/custom_errors"
	"net/http"
)

func GetStatusCodeForError(err error) int {

	switch {
	case errors.As(err, &customerrors.ValidationErr) ||
		errors.As(err, &customerrors.CredentialErr) ||
		errors.As(err, &customerrors.NotFoundOrLeaderErr):
		return http.StatusNotAcceptable
	case errors.As(err, &customerrors.ServerErr):
		return http.StatusInternalServerError
	case errors.As(err, &customerrors.PermissionErr) || errors.As(err, &customerrors.AuthErr):
		return http.StatusUnauthorized
	case errors.As(err, &customerrors.NotFoundErr):
		return http.StatusNotFound

	default:
		return http.StatusOK
	}
}
