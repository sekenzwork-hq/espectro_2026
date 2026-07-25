package handlers

import (
	"espectro/entity"
	"espectro/pkg"
	"espectro/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StaffHandlers struct {
	usecases usecases.StaffUsecases
}

func NewStaffHandlers(usecases usecases.StaffUsecases) StaffHandlers {
	return StaffHandlers{usecases: usecases}
}

func (s StaffHandlers) CreateStaff(ctx *gin.Context) {
	var entity entity.StaffFromJsonEntity
	canGo := pkg.ParseJson(ctx, &entity)

	if !canGo {
		return
	}

	newStaff, err := s.usecases.CreateStaff(entity)

	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Staff has been created", "staff": newStaff})
}
