package handlers

import (
	"espectro/entity"
	"espectro/pkg"
	"espectro/usecases"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SpeakerHandlers struct {
	usecases usecases.SpeakerUsecases
}

func NewSpeakerHandlers(usecases usecases.SpeakerUsecases) SpeakerHandlers {
	return SpeakerHandlers{usecases: usecases}
}

func (s SpeakerHandlers) CreateSpeaker(ctx *gin.Context) {

	form, formErr := ctx.MultipartForm()

	if formErr != nil {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": "Invalid form data"})
		return
	}

	fullname, fullnameExists := ctx.GetPostForm("fullname")
	bio, bioExists := ctx.GetPostForm("bio")
	country, countryExists := ctx.GetPostForm("country")
	phoneNumber, phoneNumberExists := ctx.GetPostForm("phone_number")
	email, emailExists := ctx.GetPostForm("email")
	isFeaturedStr, isFeaturedExists := ctx.GetPostForm("is_featured")

	var profilePic *multipart.FileHeader

	if !fullnameExists {
		ctx.JSON(http.StatusNotAcceptable, pkg.EmptyMessage("Provide fullname"))
	} else if !bioExists {
		ctx.JSON(http.StatusNotAcceptable, pkg.EmptyMessage("Provide bio"))
	} else if !countryExists {
		ctx.JSON(http.StatusNotAcceptable, pkg.EmptyMessage("Provide country"))
	} else if !phoneNumberExists {
		ctx.JSON(http.StatusNotAcceptable, pkg.EmptyMessage("Provide phone number"))
	} else if !emailExists {
		ctx.JSON(http.StatusNotAcceptable, pkg.EmptyMessage("Provide email"))
	} else if !isFeaturedExists {
		ctx.JSON(http.StatusNotAcceptable, pkg.EmptyMessage("Provide whether speaker is featured or not"))
	} else {

		if len(isFeaturedStr) > 5 {
			ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": "Feature should be true or false"})
			return
		}

		isFeatured, boolParseErr := strconv.ParseBool(isFeaturedStr)

		if boolParseErr != nil {
			ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": "Is feature should be true or false"})
			return
		}

		if len(form.File["profile_pic"]) >= 1 {
			profilePic = form.File["profile_pic"][0]
		}

		createdSpeaker, creationErr := s.usecases.CreateSpeaker(entity.SpeakerCreateEntity{
			Fullname:    fullname,
			ProfilePic:  profilePic,
			IsFeatured:  isFeatured,
			Bio:         bio,
			Country:     country,
			PhoneNumber: phoneNumber,
			Email:       email,
		})

		if creationErr != nil {
			code := pkg.GetStatusCodeForError(creationErr)
			ctx.JSON(code, gin.H{"status": code, "message": creationErr.Error()})
		} else {
			ctx.JSON(http.StatusCreated, gin.H{"status": 201, "message": "Speaker has been created successfully", "speaker": createdSpeaker})
		}

	}

}

func (s SpeakerHandlers) UpdateSpeaker(ctx *gin.Context) {
	form, formErr := ctx.MultipartForm()

	if formErr != nil {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": "Invalid form data"})
		return
	}

	speakerId, speakerIdExists := ctx.GetPostForm("speaker_id")

	if !speakerIdExists {
		ctx.JSON(http.StatusNotAcceptable, pkg.EmptyMessage("Provide speaker id"))
		return
	}

	fullnameForm, fullnameExists := ctx.GetPostForm("fullname")
	bioForm, bioExists := ctx.GetPostForm("bio")
	countryForm, countryExists := ctx.GetPostForm("country")
	phoneNumberForm, phoneNumberExists := ctx.GetPostForm("phone_number")
	emailForm, emailExists := ctx.GetPostForm("email")
	isFeaturedStrForm, isFeaturedExists := ctx.GetPostForm("is_featured")

	var fullname *string
	var bio *string
	var country *string
	var phoneNumber *string
	var isFeatured *bool
	var profilePic *multipart.FileHeader
	var email *string

	if fullnameExists {
		fullname = &fullnameForm
	}

	if emailExists {
		email = &emailForm
	}

	if bioExists {
		bio = &bioForm
	}

	if countryExists {
		country = &countryForm
	}

	if phoneNumberExists {
		phoneNumber = &phoneNumberForm
	}

	if isFeaturedExists {
		parsedBool, boolParseErr := strconv.ParseBool(isFeaturedStrForm)

		if boolParseErr != nil {
			ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": "Is feature should be true or false"})
			return
		}

		isFeatured = &parsedBool
	}

	if len(form.File["profile_pic"]) >= 1 {
		profilePic = form.File["profile_pic"][0]
	}

	updatedSpeaker, updationErr := s.usecases.UpdateSpeaker(speakerId, entity.SpeakerUpdateEntity{
		Fullname:    fullname,
		ProfilePic:  profilePic,
		IsFeatured:  isFeatured,
		Bio:         bio,
		Country:     country,
		PhoneNumber: phoneNumber,
		Email:       email,
	})

	if updationErr != nil {
		code := pkg.GetStatusCodeForError(updationErr)
		ctx.JSON(code, gin.H{"status": code, "message": updationErr.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Speaker has been updated successfully", "speaker": updatedSpeaker})
	}

}

func (s SpeakerHandlers) DeleteSpeaker(ctx *gin.Context) {

	speakerId, exists := ctx.GetQuery("speaker_id")

	if !exists {
		ctx.JSON(http.StatusNotAcceptable, pkg.EmptyMessage("Provide speaker id"))
		return
	}

	err := s.usecases.DeleteSpeaker(speakerId)

	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})

	} else {
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Speaker has been deleted successfully"})
	}
}

func (s SpeakerHandlers) AddSpeakerToEvent(ctx *gin.Context) {

	speakerId, speakerExists := ctx.GetQuery("speaker_id")

	if !speakerExists {
		ctx.JSON(http.StatusNotAcceptable, pkg.EmptyMessage("Provide speaker id"))
		return
	}
	eventId, eventExists := ctx.GetQuery("event_id")

	if !eventExists {
		ctx.JSON(http.StatusNotAcceptable, pkg.EmptyMessage("Provide event id"))
		return
	}

	err := s.usecases.AddSpeakerToEvent(eventId, speakerId)

	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Speaker has been added to the event successfully"})
	}
}
