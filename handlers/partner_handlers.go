package handlers

import (
	"espectro/entity"
	"espectro/pkg"
	"espectro/usecases"
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

	var entity entity.PartnerCreateEntity
	canGo := pkg.ParseJson(ctx, &entity)
	if !canGo {
		return
	}

	partner, err := p.usecases.AddPartner(entity)
	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Partner has been added", "id": partner.Id, "created_at": partner.CreatedAt})
	}

}

func (p PartnerHandlers) UpdatePartner(ctx *gin.Context) {

	var entity entity.PartnerUpdateEntity
	canGo := pkg.ParseJson(ctx, &entity)
	if !canGo {
		return
	}

	err := p.usecases.UpdatePartner(entity)
	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Partner has been updated"})

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
