package handlers

import (
	"espectro/pkg"
	"espectro/usecases"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PartnerHandlers struct {
	usecases usecases.PartnerUsecases
}

func NewPartnerHandlers(usecases usecases.PartnerUsecases) PartnerHandlers {
	return PartnerHandlers{usecases: usecases}
}

func (p PartnerHandlers) CreatePartner(ctx *gin.Context) {
	name, nameExists := ctx.GetPostForm("name")
	logo, err := ctx.FormFile("logo")

	if !nameExists {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": "Provide name"})
		return
	}

	partner, err := p.usecases.CreatePartner(name, logo)
	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
	} else {
		ctx.JSON(http.StatusCreated, gin.H{"status": 201, "message": "Partner has been created", "partner": partner})
	}

}

func (p PartnerHandlers) UpdatePartner(ctx *gin.Context) {

	var name *string
	var logo *multipart.FileHeader

	nameForm, nameExists := ctx.GetPostForm("name")
	logoForm, _ := ctx.FormFile("logo")
	partnerId, idExists := ctx.GetPostForm("partner_id")

	if !idExists {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": "Provide partner id"})
		return
	}
	if nameExists {
		name = &nameForm
	}
	if logoForm != nil {
		logo = logoForm
	}
	newPartner, err := p.usecases.UpdatePartner(partnerId, name, logo)
	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Partner has been updated", "partner": newPartner})

}

func (p PartnerHandlers) DeletePartner(ctx *gin.Context) {

	partnerId, exists := ctx.GetQuery("partner_id")

	if !exists {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": "Provide partner id"})
		return
	}

	err := p.usecases.DeletePartner(partnerId)
	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Partner has been deleted"})

}

func (p PartnerHandlers) RetrievePartner(ctx *gin.Context) {

	limit, page, limitErr := pkg.ParsePageAndLimit(ctx.Query("limit"), ctx.Query("page"))
	if limitErr != nil {
		code := pkg.GetStatusCodeForError(limitErr)
		ctx.JSON(code, gin.H{"status": code, "message": limitErr.Error()})
	}

	partners, err := p.usecases.RetrievePartner(limit, page)
	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Request was successful", "partners": partners})
}
