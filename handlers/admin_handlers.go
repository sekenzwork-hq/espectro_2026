package handlers

import (
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

	if err != nil {

		ctx.JSON(pkg.GetStatusCodeForError(err), gin.H{"status": pkg.GetStatusCodeForError(err), "message": err.Error()})
	} else {

		ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Login credentials are correct", "access_token": jwtToken})
	}

}

func (a AdminHandlers) CreateNewAdmin(ctx *gin.Context) {

	var newAdmin entity.AdminEntity

	canGo := pkg.ParseJson(ctx, &newAdmin)

	if !canGo {
		return
	}

	currentAdminId := ctx.GetString("admin_id")

	newAdminId, creationErr := a.usecases.CreateNewAdmin(newAdmin, currentAdminId)

	if creationErr != nil {
		ctx.JSON(pkg.GetStatusCodeForError(creationErr), gin.H{"status": pkg.GetStatusCodeForError(creationErr), "message": creationErr.Error()})

	} else {
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "New admin has been created", "id": newAdminId})
	}
}

func (a AdminHandlers) DeleteMemberOrVolunteer(ctx *gin.Context) {

	currentAdminId := ctx.GetString("admin_id")
	oneToDelete := ctx.Query("admin_id")

	deleted, err := a.usecases.DeleteMemberOrVolunteer(oneToDelete, currentAdminId)

	if err != nil {

		ctx.JSON(pkg.GetStatusCodeForError(err), gin.H{"status": pkg.GetStatusCodeForError(err), "message": err.Error()})

	} else if !deleted {

		ctx.JSON(http.StatusNotFound, gin.H{"status": 404, "message": "Deletion operation doesn't work whether admin doesn't exist or admin is a leader"})

	} else {

		ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Admin has been deleted"})
	}
}
