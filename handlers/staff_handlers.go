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

func (s StaffHandlers) UpdateStaff(ctx *gin.Context) {

	var entity entity.StaffUpdateEntity
	canGo := pkg.ParseJson(ctx, &entity)
	if !canGo {
		return
	}

	newStaff, err := s.usecases.UpdateStaff(entity)
	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Staff has been updated", "staff": newStaff})
}

func (s StaffHandlers) DeleteStaff(ctx *gin.Context) {

	staffId, exists := ctx.GetQuery("staff_id")
	if !exists {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": "Provide staff id"})
		return
	}

	err := s.usecases.DeleteStaff(staffId)

	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Staff has been deleted"})
}

func (s StaffHandlers) RetrieveStaffs(ctx *gin.Context) {

	limit, page, limitErr := pkg.ParsePageAndLimit(ctx.Query("limit"), ctx.Query("page"))
	if limitErr != nil {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": limitErr.Error()})
		return
	}

	staffs, err := s.usecases.RetrieveStaffs(limit, page)
	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Request was successful", "staffs": staffs})
}
