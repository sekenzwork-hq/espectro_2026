package handlers

import (
	"espectro/enums"
	"espectro/pkg"
	"espectro/usecases"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SponsorHandlers struct {
	usecases usecases.SponsorUsecases
}

func NewSponsorHandlers(usecases usecases.SponsorUsecases) SponsorHandlers {
	return SponsorHandlers{usecases: usecases}
}

func (s SponsorHandlers) CreateSponsor(ctx *gin.Context) {
	name := ctx.PostForm("name")
	amount := ctx.PostForm("amount")
	profileOrOrgImage, _ := ctx.FormFile("profile_or_org")
	sponsoredType := ctx.PostForm("sponsored_type")
	eventIds := ctx.PostFormArray("event_id")

	var amountFloat *float32

	if len(amount) != 0 {
		value, err := strconv.ParseFloat(amount, 32)
		if err != nil {
			f := float32(value)
			amountFloat = &f
		}
	}

	newSponsor, err := s.usecases.CreateSponsor(name,
		amountFloat,
		profileOrOrgImage,
		enums.SponsoredType(sponsoredType),
		eventIds,
	)

	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Sponsor has been created", "sponsor": newSponsor})
	}
}
