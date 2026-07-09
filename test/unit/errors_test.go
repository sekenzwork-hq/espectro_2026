package unit

import (
	customerrors "espectro/custom_errors"
	"espectro/pkg"
	"net/http"
	"testing"
)

func TestHttpCodeForErrors(t *testing.T) {

	serverErr := &customerrors.ServerError{OrgError: "Some error"}
	validationErr := &customerrors.ValidationError{OrgError: "Some error"}

	serverErrCode := pkg.GetStatusCodeForError(serverErr)
	validationErrCode := pkg.GetStatusCodeForError(validationErr)

	if serverErrCode != http.StatusInternalServerError {
		t.Errorf("Status code generation for error failed. Error : %v, Code : %v", serverErr, serverErrCode)
	}
	if validationErrCode != http.StatusNotAcceptable {
		t.Errorf("Status code generation for error failed. Error : %v, Code : %v", validationErr, validationErrCode)
	}

}
