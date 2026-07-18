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

func (s SponsorHandlers) UpdateSponsor(ctx *gin.Context) {
	sponsorId := ctx.PostForm("sponsor_id")
	nameForm := ctx.PostForm("name")
	amountForm := ctx.PostForm("amount")
	profileOrOrgImage, _ := ctx.FormFile("profile_or_org")
	sponsoredTypeForm := ctx.PostForm("sponsored_type")

	var name *string
	var amountFloat *float32
	var sponsoredType *enums.SponsoredType

	if len(nameForm) != 0 {
		name = &nameForm
	}

	if len(amountForm) != 0 {
		value, err := strconv.ParseFloat(amountForm, 32)

		if err != nil {
			ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": "Invalid amount"})
			return
		}
		f := float32(value)
		amountFloat = &f

	}

	if len(sponsoredTypeForm) != 0 {
		st := enums.SponsoredType(sponsoredTypeForm)
		sponsoredType = &st
	}

	newSponsor, err := s.usecases.UpdateSponsor(sponsorId, name, amountFloat, profileOrOrgImage, sponsoredType)
	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Sponsor has been updated", "sponsor": newSponsor})

}

func (s SponsorHandlers) DeleteSponsor(ctx *gin.Context) {

	sponsorId := ctx.Query("sponsor_id")

	err := s.usecases.DeleteSponsor(sponsorId)
	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Sponsor has been deleted"})
}

func (s SponsorHandlers) RetrieveSponsors(ctx *gin.Context) {

	limit, page, limitErr := pkg.ParsePageAndLimit(ctx.Query("limit"), ctx.Query("page"))
	if limitErr != nil {
		code := pkg.GetStatusCodeForError(limitErr)
		ctx.JSON(code, gin.H{"status": code, "message": limitErr.Error()})
		return
	}

	sponsors, err := s.usecases.RetrieveSponsors(limit, page)
	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Request was successful", "sponsors": sponsors})

}
