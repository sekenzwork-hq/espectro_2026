package handlers

import (
	"espectro/enums"
	"espectro/pkg"
	"espectro/usecases"
	"fmt"
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
	imageForm := form.File["image"]

	fmt.Println("Images from controller : ", imageForm)

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

func (s SpectrumHandlers) UpdateSpectrum(ctx *gin.Context) {

	form, err := ctx.MultipartForm()

	if err != nil {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": "Invalid form"})
		return
	}

	spectrumIdForm := form.Value["spectrum_id"]

	if len(spectrumIdForm) == 0 {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": "Provide spectrum id"})
		return
	}
	nameForm := form.Value["name"]
	shortDescriptionForm := form.Value["short_description"]
	descriptionForm := form.Value["description"]
	statusForm := form.Value["status"]
	logoFileForm := form.File["logo"]
	videoFileForm := form.File["video"]
	imageForm := form.File["image"]

	var name *string
	var shortDes *string
	var des *string
	var status *string
	var logoFile *multipart.FileHeader
	var videoFile *multipart.FileHeader

	if len(nameForm) != 0 {
		name = &nameForm[0]
	}
	if len(shortDescriptionForm) != 0 {
		shortDes = &shortDescriptionForm[0]
	}
	if len(descriptionForm) != 0 {
		des = &descriptionForm[0]
	}
	if len(statusForm) != 0 {
		status = &statusForm[0]
	}
	if len(logoFileForm) != 0 {
		logoFile = logoFileForm[0]
	}
	if len(videoFileForm) != 0 {
		videoFile = videoFileForm[0]
	}

	updatedMedia, err := s.usecases.UpdateSpectrum(
		spectrumIdForm[0],
		name,
		shortDes,
		des, (*enums.SpectrumStatus)(status),
		imageForm,
		videoFile,
		logoFile,
	)

	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
	} else {

		json := map[string]any{}

		if updatedMedia.LogoUrl != nil {
			json["logo_url"] = updatedMedia.LogoUrl
		}
		if updatedMedia.VideoUrl != nil {
			json["video_url"] = updatedMedia.VideoUrl
		}
		if len(updatedMedia.ImagesUrls) != 0 {
			json["image_urls"] = updatedMedia.ImagesUrls
		}

		if len(json) != 0 {
			ctx.JSON(http.StatusAccepted, gin.H{"status": 202, "message": "Spectrum has been updated", "updated_media": json})
		} else {
			ctx.JSON(http.StatusAccepted, gin.H{"status": 202, "message": "Spectrum has been updated"})
		}
	}
}

func (s SpectrumHandlers) DeleteSpectrum(ctx *gin.Context) {

	spectrumId := ctx.Query("spectrum_id")

	err := s.usecases.DeleteSpectrum(spectrumId)

	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Spectrum has been deleted"})

}
