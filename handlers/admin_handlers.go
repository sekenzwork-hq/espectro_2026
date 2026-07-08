package handlers

import (
	"errors"
	customerrors "espectro/custom_errors"
	"espectro/entity"
	"espectro/pkg"
	"espectro/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdminHandlers struct {
	usecases usecases.AdminUsecases
}

func NewAdminHandlers(usecases usecases.AdminUsecases) AdminHandlers {
	return AdminHandlers{
		usecases: usecases,
	}
}

func (a AdminHandlers) Login(ctx *gin.Context) {

	var enteredData entity.AdminEnteredLoginCredentials
	parsed := pkg.ParseJson(ctx, &enteredData)

	if !parsed {
		return
	}

	jwtToken, err := a.usecases.Login(enteredData.Email, enteredData.Password)

	if _, ok := errors.AsType[*customerrors.CredentialsError](err); ok {

		ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": err.Error()})

	} else if _, ok := errors.AsType[*customerrors.ServerError](err); ok {

		ctx.JSON(http.StatusInternalServerError, gin.H{"status": 500, "message": err.Error()})

	} else {

		ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Login credentials are correct", "access_token": jwtToken})
	}

}
