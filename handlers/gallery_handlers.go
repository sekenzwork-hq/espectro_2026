package handlers

import (
	"espectro/pkg"
	"espectro/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GalleryHandlers struct {
	usecases usecases.GalleryUsecases
}

func NewGalleryHandlers(usecases usecases.GalleryUsecases) GalleryHandlers {
	return GalleryHandlers{usecases: usecases}
}

func (g GalleryHandlers) CreateGallery(ctx *gin.Context) {

	name, nameExists := ctx.GetPostForm("name")
	venueId, venueIdExists := ctx.GetPostForm("venue_id")
	imagesForms, _ := ctx.MultipartForm()
	images := imagesForms.File["images"]

	if !nameExists {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": "Provide name for gallery"})
		return
	}

	if !venueIdExists {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": "Provide venue id"})
		return
	}

	gallery, err := g.usecases.CreateGallery(name, venueId, images)

	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": 201, "message": "Gallery has been created", "gallery": gallery})
}
