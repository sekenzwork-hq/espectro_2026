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

func (p PartnerHandlers) AddPartner(ctx *gin.Context) {
	name := ctx.PostForm("name")
	logo, err := ctx.FormFile("logo")

	partner, err := p.usecases.AddPartner(name, logo)
	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Partner has been added", "partner": partner})
	}

}

func (p PartnerHandlers) UpdatePartner(ctx *gin.Context) {

	var name *string
	var logo *multipart.FileHeader

	nameForm := ctx.PostForm("name")
	logoForm, _ := ctx.FormFile("logo")
	partnerId := ctx.PostForm("partner_id")

	if len(nameForm) != 0 {
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

	partnerId := ctx.Query("partner_id")

	err := p.usecases.DeletePartner(partnerId)

	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Partner has been deleted"})

}

func (p PartnerHandlers) RetrievePartner(ctx *gin.Context) {

	limit, page := pkg.ParsePageAndLimit(ctx.Query("limit"), ctx.Query("page"))

	partners, err := p.usecases.RetrievePartner(limit, page)

	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Request was successful", "partners": partners})
}
