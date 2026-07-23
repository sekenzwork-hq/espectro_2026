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

	name, nameExists := ctx.GetQuery("name")
	shortDescription, shortDesExists := ctx.GetQuery("short_description")
	description, desExists := ctx.GetQuery("description")
	status, statusExists := ctx.GetQuery("status")
	form, _ := ctx.MultipartForm()
	logoFileForm := form.File["logo"]
	videoFileForm := form.File["video"]
	imageForm := form.File["images"]

	if !nameExists {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": "Provide name"})
		return
	}

	if !shortDesExists {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": "Provide short description"})
		return
	}
	if !desExists {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": "Provide description"})
		return
	}
	if !statusExists {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": "Provide status"})
		return
	}

	var video *multipart.FileHeader
	var logo *multipart.FileHeader

	if len(videoFileForm) != 0 {
		video = videoFileForm[0]
	}
	if len(logoFileForm) != 0 {
		logo = logoFileForm[0]
	}
	spectrum, err := s.usecases.CreateSpectrum(name, shortDescription, description, enums.SpectrumStatus(status), imageForm, video, logo)
	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
	} else {
		ctx.JSON(http.StatusCreated, gin.H{"status": 201, "message": "Spectrum has been created", "spectrum": spectrum})
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

	newSpectrum, err := s.usecases.UpdateSpectrum(
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
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Spectrum has been updated", "spectrum": newSpectrum})

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

func (s SpectrumHandlers) RetrieveSpectrums(ctx *gin.Context) {

	limit, page, limitErr := pkg.ParsePageAndLimit(ctx.Query("limit"), ctx.Query("page"))

	if limitErr != nil {
		code := pkg.GetStatusCodeForError(limitErr)
		ctx.JSON(code, gin.H{"status": code, "message": limitErr.Error()})
	}

	spectrums, spectrumsErr := s.usecases.RetrieveSpectrums(limit, page)

	if spectrumsErr != nil {
		code := pkg.GetStatusCodeForError(spectrumsErr)
		ctx.JSON(code, gin.H{"status": code, "message": spectrumsErr.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Request was successful", "spectrums": spectrums})

}
