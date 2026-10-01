package handlers

import (
	"espectro/entity"
	"espectro/enums"
	"espectro/pkg"
	"espectro/usecases"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type CompetitionHandlers struct {
	usecases usecases.CompetitionUsecases
}

func NewCompetitionHandlers(usecases usecases.CompetitionUsecases) CompetitionHandlers {
	return CompetitionHandlers{usecases: usecases}
}

func (c CompetitionHandlers) CreateCompetition(ctx *gin.Context) {

	form, formErr := ctx.MultipartForm()

	if formErr != nil {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": "Invalid form data"})
		return
	}

	title, titleExists := ctx.GetPostForm("title")
	eventId, eventIdExists := ctx.GetPostForm("event_id")
	des, desExists := ctx.GetPostForm("description")
	rules, rulesExists := ctx.GetPostForm("rules")
	status, statusExists := ctx.GetPostForm("status")
	oTimeForm, oTimeExists := ctx.GetPostForm("opening_time")
	cTimeForm, cTimeExists := ctx.GetPostForm("closing_time")
	rTimeForm, rTimeExists := ctx.GetPostForm("result_time")
	regFeeForm, regFeeExists := ctx.GetPostForm("registration_fee")
	prizeId, prizeIdExists := ctx.GetPostForm("prize_id")
	venueId, venueIdExists := ctx.GetPostForm("venue_id")

	if !titleExists {
		ctx.JSON(http.StatusNotAcceptable, pkg.EmptyMessage("Provide title"))
	} else if !eventIdExists {
		ctx.JSON(http.StatusNotAcceptable, pkg.EmptyMessage("Provide event id"))
	} else if !desExists {
		ctx.JSON(http.StatusNotAcceptable, pkg.EmptyMessage("Provide description"))
	} else if !rulesExists {
		ctx.JSON(http.StatusNotAcceptable, pkg.EmptyMessage("Provide rules"))
	} else if !statusExists {
		ctx.JSON(http.StatusNotAcceptable, pkg.EmptyMessage("Provide status"))
	} else if !oTimeExists {
		ctx.JSON(http.StatusNotAcceptable, pkg.EmptyMessage("Provide opening time"))
	} else if !cTimeExists {
		ctx.JSON(http.StatusNotAcceptable, pkg.EmptyMessage("Provide closing time"))
	} else if !rTimeExists {
		ctx.JSON(http.StatusNotAcceptable, pkg.EmptyMessage("Provide result time"))
	} else if !regFeeExists {
		ctx.JSON(http.StatusNotAcceptable, pkg.EmptyMessage("Provide registration fee"))
	} else if !prizeIdExists {
		ctx.JSON(http.StatusNotAcceptable, pkg.EmptyMessage("Provide prize id"))
	} else if !venueIdExists {
		ctx.JSON(http.StatusNotAcceptable, pkg.EmptyMessage("Provide venue id"))
	} else {

		var regFee float32
		var logo *multipart.FileHeader

		fValue, err := strconv.ParseFloat(regFeeForm, 32)

		if err != nil {
			ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": "Invalid registration fee"})
			return
		}

		regFee = float32(fValue)

		openingTime, oTimeParseErr := time.Parse(time.RFC3339, oTimeForm)

		if oTimeParseErr != nil {
			ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": "Invalid opening time format"})
			return
		}

		closingTime, cTimeParseErr := time.Parse(time.RFC3339, cTimeForm)

		if cTimeParseErr != nil {
			ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": "Invalid closing time"})
			return
		}
		resultTime, rTimeParseErr := time.Parse(time.RFC3339, rTimeForm)

		if rTimeParseErr != nil {
			ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": "Invalid result time"})
			return
		}

		if len(form.File["logo"]) >= 1 {
			logo = form.File["logo"][0]
		}

		createdCompetition, creationErr := c.usecases.CreateCompetition(entity.CompetitionCreateEntity{
			EventId:         eventId,
			Title:           title,
			Description:     des,
			Rules:           rules,
			Status:          enums.CompetitionStatusEnum(status),
			OpeningTime:     openingTime,
			ClosingTime:     closingTime,
			ResultTime:      resultTime,
			RegistrationFee: regFee,
			PrizeId:         prizeId,
			VenueId:         venueId,
			Logo:            logo,
		})

		if creationErr != nil {
			code := pkg.GetStatusCodeForError(creationErr)
			ctx.JSON(code, gin.H{"status": code, "message": creationErr.Error()})
		} else {
			ctx.JSON(http.StatusCreated, gin.H{"status": 201, "message": "Competition has been created successfully", "competition": createdCompetition})
		}
	}
}

func (c CompetitionHandlers) UpdateCompetition(ctx *gin.Context) {

	form, formErr := ctx.MultipartForm()

	if formErr != nil {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": "Invalid form data"})
		return
	}

	competitionId, exists := ctx.GetPostForm("competition_id")

	if !exists {
		ctx.JSON(http.StatusNotAcceptable, pkg.EmptyMessage("Provide competition id"))
		return
	}

	titleForm, titleExists := ctx.GetPostForm("title")
	eventIdForm, eventIdExists := ctx.GetPostForm("event_id")
	desForm, desExists := ctx.GetPostForm("description")
	rulesForm, rulesExists := ctx.GetPostForm("rules")
	statusForm, statusExists := ctx.GetPostForm("status")
	oTimeForm, oTimeExists := ctx.GetPostForm("opening_time")
	cTimeForm, cTimeExists := ctx.GetPostForm("closing_time")
	rTimeForm, rTimeExists := ctx.GetPostForm("result_time")
	regFeeForm, regFeeExists := ctx.GetPostForm("registration_fee")
	prizeIdForm, prizeIdExists := ctx.GetPostForm("prize_id")
	venueIdForm, venueIdExists := ctx.GetPostForm("venue_id")

	var (
		title           *string
		eventId         *string
		description     *string
		rules           *string
		status          *enums.CompetitionStatusEnum
		openingTime     *time.Time
		closingTime     *time.Time
		resultTime      *time.Time
		registrationFee *float32
		prizeId         *string
		venueId         *string
		logo            *multipart.FileHeader
	)

	if titleExists {
		title = &titleForm
	}

	if desExists {
		description = &desForm
	}

	if rulesExists {
		rules = &rulesForm
	}

	if prizeIdExists {
		prizeId = &prizeIdForm
	}

	if eventIdExists {
		eventId = &eventIdForm
	}

	if venueIdExists {
		venueId = &venueIdForm
	}

	if oTimeExists {
		oTime, oTimeParseErr := time.Parse(time.RFC3339, oTimeForm)
		if oTimeParseErr != nil {
			ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": "Invalid opening time format"})
			return
		}
		openingTime = &oTime
	}

	if cTimeExists {
		cTime, cTimeParseErr := time.Parse(time.RFC3339, cTimeForm)
		if cTimeParseErr != nil {
			ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": "Invalid closing time format"})
			return
		}
		closingTime = &cTime
	}

	if rTimeExists {
		rTime, rTimeParseErr := time.Parse(time.RFC3339, rTimeForm)
		if rTimeParseErr != nil {
			ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": "Invalid opening time format"})
			return
		}
		resultTime = &rTime
	}

	if regFeeExists {

		fValue, parseErr := strconv.ParseFloat(regFeeForm, 32)

		if parseErr != nil {
			ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": "Invalid registration fee"})
			return
		}
		f := float32(fValue)
		registrationFee = &f
	}

	if statusExists {
		e := enums.CompetitionStatusEnum(statusForm)
		status = &e
	}

	if len(form.File["logo"]) >= 1 {
		logo = form.File["logo"][0]
	}

	updatedCompetition, updationErr := c.usecases.UpdateCompetition(competitionId, entity.CompetitionUpdateEntity{
		EventId:         eventId,
		Title:           title,
		Description:     description,
		Rules:           rules,
		Status:          status,
		OpeningTime:     openingTime,
		ClosingTime:     closingTime,
		ResultTime:      resultTime,
		RegistrationFee: registrationFee,
		PrizeId:         prizeId,
		VenueId:         venueId,
		Logo:            logo,
	})

	if updationErr != nil {
		code := pkg.GetStatusCodeForError(updationErr)
		ctx.JSON(code, gin.H{"status": code, "message": updationErr.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Competition has been updated successfully", "competition": updatedCompetition})
	}
}

func (c CompetitionHandlers) DeleteCompetition(ctx *gin.Context) {

	competitionId, exists := ctx.GetQuery("competition_id")

	if !exists {
		ctx.JSON(http.StatusNotAcceptable, pkg.EmptyMessage("Provide competition id"))
		return
	}

	err := c.usecases.DeleteCompetition(competitionId)

	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Competition has been deleted successfully"})
	}
}

func (c CompetitionHandlers) RetrieveCompetitions(ctx *gin.Context) {

	eventId, exists := ctx.GetQuery("event_id")

	if !exists {
		ctx.JSON(http.StatusNotAcceptable, pkg.EmptyMessage("Provide event id"))
		return
	}

	page := pkg.ParsePageSetMin1(ctx.Query("page"))

	competitions, err := c.usecases.RetrieveCompetitions(eventId, page)

	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Request was successful", "competitions": competitions})
	}

}
