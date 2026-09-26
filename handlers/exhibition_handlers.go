package handlers

import (
	"espectro/entity"
	"espectro/pkg"
	"espectro/usecases"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ExhibitionHandlers struct {
	exhibitonUsecases usecases.ExhibitionUsecases
}

func NewExhibitionHandlers(exhibitionUsecases usecases.ExhibitionUsecases) ExhibitionHandlers {
	return ExhibitionHandlers{exhibitonUsecases: exhibitionUsecases}
}

func (e ExhibitionHandlers) CreateExhibition(ctx *gin.Context) {

	form, formError := ctx.MultipartForm()

	if formError != nil {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": "Invalid form data"})
		return
	}

	title, titleExists := ctx.GetPostForm("title")
	desc, descExists := ctx.GetPostForm("description")
	eventId, eventIdExists := ctx.GetPostForm("event_id")
	orgId, orgIdExists := ctx.GetPostForm("organization_id")
	category, categoryExists := ctx.GetPostForm("category")
	userId, userIdExists := ctx.GetPostForm("user_id")

	images := form.File["images"]

	if !userIdExists {
		ctx.JSON(http.StatusNotAcceptable, pkg.EmptyMessage("Provide user id"))
	} else if !titleExists {
		ctx.JSON(http.StatusNotAcceptable, pkg.EmptyMessage("Provide title"))
	} else if !descExists {
		ctx.JSON(http.StatusNotAcceptable, pkg.EmptyMessage("Provide description"))
	} else if !eventIdExists {
		ctx.JSON(http.StatusNotAcceptable, pkg.EmptyMessage("Provide event id"))
	} else if !orgIdExists {
		ctx.JSON(http.StatusNotAcceptable, pkg.EmptyMessage("Provide organization id"))
	} else if !categoryExists {
		ctx.JSON(http.StatusNotAcceptable, pkg.EmptyMessage("Provide category"))
	} else {

		createdExhibition, creationError := e.exhibitonUsecases.CreateExhibition(entity.ExhibitionRawEntity{
			EventId:         eventId,
			UserId:          &userId,
			Category:        category,
			OrganizationId:  orgId,
			ItemTitle:       title,
			ItemImages:      images,
			ItemDescription: desc,
		})

		if creationError != nil {
			code := pkg.GetStatusCodeForError(creationError)
			ctx.JSON(code, gin.H{"status": code, "message": creationError.Error()})
		} else {
			ctx.JSON(http.StatusCreated, gin.H{"status": 201, "message": "Exhibtion has been created", "exhibtion": createdExhibition})
		}
	}
}

func (e ExhibitionHandlers) UpdateExhibitionFromUserSide(ctx *gin.Context) {

	form, formErr := ctx.MultipartForm()

	if formErr != nil {
		ctx.JSON(http.StatusNotAcceptable, pkg.EmptyMessage("Invalid form data"))
		return
	}

	exhibitionId, exhibitionIdExists := ctx.GetPostForm("exhibition_id")

	if !exhibitionIdExists {
		ctx.JSON(http.StatusNotAcceptable, pkg.EmptyMessage("Provide exhibition id"))
		return
	}

	titleForm, titleExists := ctx.GetPostForm("title")
	descForm, descExists := ctx.GetPostForm("description")
	eventIdForm, eventIdExists := ctx.GetPostForm("event_id")
	orgIdForm, orgIdExists := ctx.GetPostForm("organization_id")
	categoryForm, categoryExists := ctx.GetPostForm("category")
	userId, userIdExists := ctx.Get("user_id")

	if !userIdExists {
		ctx.JSON(http.StatusNotAcceptable, pkg.EmptyMessage("Provide user id"))
		return
	}

	images := form.File["images"]

	var title *string
	var desc *string
	var eventId *string
	var orgId *string
	var category *string

	if titleExists {
		title = &titleForm
	}
	if descExists {
		desc = &descForm
	}
	if eventIdExists {
		eventId = &eventIdForm
	}

	if orgIdExists {
		orgId = &orgIdForm
	}
	if categoryExists {
		category = &categoryForm
	}

	userIdStr := fmt.Sprint(userId)

	updatedExhibition, updationError := e.exhibitonUsecases.UpdateExhibitionFromUserSide(exhibitionId, entity.ExhibitionRawUpdateEntity{
		EventId:         eventId,
		UserId:          &userIdStr,
		Category:        category,
		OrganizationId:  orgId,
		ItemTitle:       title,
		ItemImages:      images,
		ItemDescription: desc,
	})

	if updationError != nil {
		code := pkg.GetStatusCodeForError(updationError)
		ctx.JSON(code, gin.H{"status": code, "message": updationError.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Exhibition has been updated successfully", "exhibition": updatedExhibition})
	}
}

func (e ExhibitionHandlers) UpdateExhibitionFromAdminSide(ctx *gin.Context) {

	var exhibitionDataToUpdate entity.ExhibitionUpdateAdminEntity

	parsed := pkg.ParseJson(ctx, &exhibitionDataToUpdate)

	if !parsed {
		return
	}

	updatedExhibition, updationError := e.exhibitonUsecases.UpdateExhibitionFromAdminSide(exhibitionDataToUpdate.Id, entity.ExhibitionRawUpdateEntity{
		TokenNumber:     exhibitionDataToUpdate.TokenNumber,
		BoothNumber:     exhibitionDataToUpdate.BoothNumber,
		AvailableSqft:   exhibitionDataToUpdate.AvailableSqft,
		AssignedStaffId: exhibitionDataToUpdate.AssignedStaffId,
		ApprovedBy:      exhibitionDataToUpdate.ApprovedBy,
		Status:          exhibitionDataToUpdate.Status,
	})

	if updationError != nil {
		code := pkg.GetStatusCodeForError(updationError)
		ctx.JSON(code, gin.H{"status": code, "message": updationError.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Exhibition has been updated successfully", "exhibition": updatedExhibition})
	}
}

func (e ExhibitionHandlers) DeleteExhibition(ctx *gin.Context) {

	exhibitionId, exists := ctx.GetQuery("exhibition_id")

	if !exists {
		ctx.JSON(http.StatusNotAcceptable, pkg.EmptyMessage("Provide exhibition id"))
		return
	}

	deletionErr := e.exhibitonUsecases.DeleteExhibition(exhibitionId)

	if deletionErr != nil {
		code := pkg.GetStatusCodeForError(deletionErr)
		ctx.JSON(code, gin.H{"status": code, "message": deletionErr.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Exhibition has been deleted"})
}

func (e ExhibitionHandlers) RetrieveExhibitionFromUserSide(ctx *gin.Context) {

	userId := ctx.GetString("user_id")
	page, _ := strconv.Atoi(ctx.Query("page"))

	if page <= 0 {
		page = 1
	}

	exhibitions, err := e.exhibitonUsecases.RetrieveExhibitionFromUserSide(userId, page)

	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Request was successful", "exhibitions": exhibitions})
	}

}

func (e ExhibitionHandlers) RetrieveExhibitionFromAdminSide(ctx *gin.Context) {

	page, _ := strconv.Atoi(ctx.Query("page"))

	if page <= 0 {
		page = 1
	}

	exhibitions, err := e.exhibitonUsecases.RetrieveExhibitionFromAdminSide(page)

	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Request was successful", "exhibitions": exhibitions})
	}
}
