package handlers

import (
	"espectro/entity"
	"espectro/pkg"
	"espectro/usecases"
	"net/http"

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

	images := form.File["images"]

	if !titleExists {
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

		createdExhibition, creationError := e.exhibitonUsecases.CreateExhibition(entity.ExhibitionCreateEntity{
			EventId:         eventId,
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
