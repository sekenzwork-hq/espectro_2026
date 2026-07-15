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

	var entity entity.PartnerFromJsonEntity
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
