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
	imagesForms, _ := ctx.MultipartForm()
	images := imagesForms.File["images"]

	if !nameExists {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": "Provide name for gallery"})
		return
	}

	gallery, err := g.usecases.CreateGallery(name, images)
	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": 201, "message": "Gallery has been created", "gallery": gallery})
}

func (g GalleryHandlers) UpdateGallery(ctx *gin.Context) {

	galleryId, galleryIdExists := ctx.GetPostForm("gallery_id")

	if !galleryIdExists {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": "Provide gallery id"})
		return
	}
	nameForm, nameExists := ctx.GetPostForm("name")
	imagesForm, _ := ctx.MultipartForm()

	images := imagesForm.File["images"]
	var name *string

	if nameExists {
		name = &nameForm
	}

	newGallery, updationErr := g.usecases.UpdateGallery(galleryId, name, images)

	if updationErr != nil {
		code := pkg.GetStatusCodeForError(updationErr)
		ctx.JSON(code, gin.H{"status": code, "message": updationErr.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Gallery has been updated", "gallery": newGallery})

}

func (g GalleryHandlers) DeleteGallery(ctx *gin.Context) {

	galleryId, exists := ctx.GetQuery("gallery_id")
	if !exists {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": "Provide gallery id"})
		return
	}

	err := g.usecases.DeleteGallery(galleryId)
	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Gallery has been deleted"})
}

func (g GalleryHandlers) RetrieveGalleries(ctx *gin.Context) {

	limit, page, limitErr := pkg.ParsePageAndLimit(ctx.Query("limit"), ctx.Query("page"))
	if limitErr != nil {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": limitErr.Error()})
		return
	}

	galleries, err := g.usecases.RetrieveGalleries(limit, page)

	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Request was successful", "galleries": galleries})
}
