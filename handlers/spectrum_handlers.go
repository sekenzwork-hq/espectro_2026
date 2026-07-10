package handlers

import (
	"espectro/enums"
	"espectro/pkg"
	"espectro/usecases"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SpectrumHandlers struct {
	usecases usecases.SpectrumUsecases
}

func NewSpectrumHandlers(usecases usecases.SpectrumUsecases) SpectrumHandlers {
	return SpectrumHandlers{usecases: usecases}
}

func (s SpectrumHandlers) CreateSpectrum(ctx *gin.Context) {

	form, err := ctx.MultipartForm()

	if err != nil {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": "Invalid form"})
		return
	}

	nameForm := form.Value["name"]
	shortDescriptionForm := form.Value["short_description"]
	descriptionForm := form.Value["description"]
	statusForm := form.Value["status"]
	logoFileForm := form.File["logo"]
	videoFileForm := form.File["video"]
	imageForm := form.File["images"]

	if len(nameForm) == 0 {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": "Invalid name"})
	} else if len(shortDescriptionForm) == 0 {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": "Invalid short description"})
	} else if len(descriptionForm) == 0 {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": "Invalid description"})
	} else if len(statusForm) == 0 || !enums.SpectrumStatus(statusForm[0]).IsValid() {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": "Invalid status"})
	} else if len(imageForm) > 10 {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": "Maximum number of images is 10"})
		return
	} else {
		name := nameForm[0]
		shortDescription := shortDescriptionForm[0]
		description := descriptionForm[0]
		status := enums.SpectrumStatus(statusForm[0])

		var video *multipart.FileHeader
		var logo *multipart.FileHeader

		if len(videoFileForm) != 0 {
			video = videoFileForm[0]
		}
		if len(logoFileForm) != 0 {
			logo = logoFileForm[0]
		}

		spectrum, err := s.usecases.CreateSpectrum(name, shortDescription, description, status, imageForm, video, logo)
		if err != nil {
			code := pkg.GetStatusCodeForError(err)
			ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
		} else {
			ctx.JSON(http.StatusCreated, gin.H{"status": 201, "message": "Spectrum has been created", "spectrum": spectrum})
		}

	}

}
