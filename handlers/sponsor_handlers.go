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
	name, nameExists := ctx.GetPostForm("name")
	amount, amountExists := ctx.GetPostForm("amount")
	profileOrOrgImage, _ := ctx.FormFile("profile_or_org")
	sponsoredType := ctx.PostForm("sponsored_type")
	eventIds := ctx.PostFormArray("event_id")

	if !nameExists {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": "Provide name"})
		return
	}

	var amountFloat *float32
	if amountExists {
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
		ctx.JSON(http.StatusCreated, gin.H{"status": 201, "message": "Sponsor has been created", "sponsor": newSponsor})
	}
}

func (s SponsorHandlers) UpdateSponsor(ctx *gin.Context) {
	sponsorId, sponsorIdExists := ctx.GetPostForm("sponsor_id")
	nameForm, nameExists := ctx.GetPostForm("name")
	amountForm, amountExists := ctx.GetPostForm("amount")
	profileOrOrgImage, _ := ctx.FormFile("profile_or_org")
	sponsoredTypeForm, typeExists := ctx.GetPostForm("sponsored_type")

	if !sponsorIdExists {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": "Provide sponsor id"})
		return
	}

	var name *string
	var amountFloat *float32
	var sponsoredType *enums.SponsoredType

	if nameExists {
		name = &nameForm
	}

	if amountExists {
		value, err := strconv.ParseFloat(amountForm, 32)
		if err != nil {
			ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": "Invalid amount"})
			return
		}
		f := float32(value)
		amountFloat = &f

	}

	if typeExists {
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

	sponsorId, exists := ctx.GetQuery("sponsor_id")

	if !exists {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": "Provide sponsor id"})
		return
	}

	err := s.usecases.DeleteSponsor(sponsorId)
	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Sponsor has been deleted"})
}

func (s SponsorHandlers) RetrieveSponsors(ctx *gin.Context) {

	var eventId *string
	eventIdQ, exists := ctx.GetQuery("event_id")
	if exists {
		eventId = &eventIdQ
	}

	limit, page, limitErr := pkg.ParsePageAndLimit(ctx.Query("limit"), ctx.Query("page"))
	if limitErr != nil {
		code := pkg.GetStatusCodeForError(limitErr)
		ctx.JSON(code, gin.H{"status": code, "message": limitErr.Error()})
		return
	}

	sponsors, err := s.usecases.RetrieveSponsors(eventId, limit, page)
	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Request was successful", "sponsors": sponsors})

}
