package handlers

import (
	"espectro/entity"
	"espectro/pkg"
	"espectro/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type VenueHandlers struct {
	usecases usecases.VenueUsecases
}

func NewVenueHandlers(usecases usecases.VenueUsecases) VenueHandlers {
	return VenueHandlers{usecases: usecases}
}

func (v VenueHandlers) CreateVenue(ctx *gin.Context) {

	var venue entity.VenueCreateEntity

	canGo := pkg.ParseJson(ctx, &venue)
	if !canGo {
		return
	}

	venueId, err := v.usecases.CreateVenue(venue)

	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": 201, "message": "Venue has been created", "id": venueId})

}

func (v VenueHandlers) DeleteVenue(ctx *gin.Context) {

	venueId := ctx.Query("venue_id")

	err := v.usecases.DeleteVenue(venueId)

	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Venue has been deleted"})
	}
}

func (v VenueHandlers) UpdateVenue(ctx *gin.Context) {

	var newVenue entity.VenueUpdateEntity

	canGo := pkg.ParseJson(ctx, &newVenue)

	if !canGo {
		return
	}

	err := v.usecases.UpdateVenue(newVenue.Id, newVenue)

	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Venue has been updated"})
	}
}

func (v VenueHandlers) RetrieveVenue(ctx *gin.Context) {

	limit, page := pkg.ParsePageAndLimit(ctx.Query("limit"), ctx.Query("page"))

	venue, err := v.usecases.RetrieveVenue(page, limit)

	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Request was successful", "venue": venue})
}
