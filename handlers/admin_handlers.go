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
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
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
		code := pkg.GetStatusCodeForError(creationErr)
		ctx.JSON(code, gin.H{"status": code, "message": creationErr.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "New admin has been created", "id": newAdminId})
	}
}

func (a AdminHandlers) DeleteMemberOrVolunteer(ctx *gin.Context) {

	currentAdminId := ctx.GetString("admin_id")
	oneToDelete := ctx.Query("admin_id")

	deleted, err := a.usecases.DeleteMemberOrVolunteer(oneToDelete, currentAdminId)

	if err != nil {

		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})

	} else if !deleted {

		ctx.JSON(http.StatusNotFound, gin.H{"status": 404, "message": "Deletion operation doesn't work whether admin doesn't exist or admin is a leader"})

	} else {

		ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Admin has been deleted"})
	}
}

func (a AdminHandlers) UpdateCurrentAdmin(ctx *gin.Context) {

	var adminEmailAndFullname entity.AdminUpdateEmailAndFullnameEntity

	canGo := pkg.ParseJson(ctx, &adminEmailAndFullname)

	if !canGo {
		return
	}

	currentAdminId := ctx.GetString("admin_id")

	err := a.usecases.UpdateCurrentAdmin(currentAdminId, adminEmailAndFullname.Email, adminEmailAndFullname.Fullname)

	if err != nil {
		statCode := pkg.GetStatusCodeForError(err)
		ctx.JSON(statCode, gin.H{"status": statCode, "message": err.Error()})
	} else {
		ctx.JSON(http.StatusAccepted, gin.H{"status": 202, "message": "Updated successfully"})
	}
}
