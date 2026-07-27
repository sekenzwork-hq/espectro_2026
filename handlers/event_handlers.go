package handlers

import (
	"espectro/entity"
	"espectro/pkg"
	"espectro/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type EventHandlers struct {
	usecases usecases.EventUsecases
}

func NewEventHandlers(usecases usecases.EventUsecases) EventHandlers {
	return EventHandlers{usecases: usecases}
}

func (e EventHandlers) CreateEvent(ctx *gin.Context) {

	var entity entity.EventCreateEntity

	canGo := pkg.ParseJson(ctx, &entity)
	if !canGo {
		return
	}

	event, err := e.usecases.CreateEvent(entity)
	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
	} else {
		ctx.JSON(http.StatusCreated, gin.H{"status": 201, "message": "Event has been created", "event": event})
	}
}

func (e EventHandlers) UpdateEvent(ctx *gin.Context) {
	var entity entity.EventUpdateEntity

	canGo := pkg.ParseJson(ctx, &entity)
	if !canGo {
		return
	}

	newEvent, err := e.usecases.UpdateEvent(entity)
	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Event has been updated", "event": newEvent})
	}

}

func (e EventHandlers) DeleteEvent(ctx *gin.Context) {

	eventId, exists := ctx.GetQuery("event_id")

	if !exists {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": "Provide event id"})
		return
	}

	err := e.usecases.DeleteEvent(eventId)
	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Event has been deleted"})
	}
}

func (e EventHandlers) RetrieveEvents(ctx *gin.Context) {

	var spectrumId *string
	spectrumIdQ, exists := ctx.GetQuery("spectrum_id")
	if exists {
		spectrumId = &spectrumIdQ
	}

	limit, page, limitErr := pkg.ParsePageAndLimit(ctx.Query("limit"), ctx.Query("page"))

	if limitErr != nil {
		code := pkg.GetStatusCodeForError(limitErr)
		ctx.JSON(code, gin.H{"status": code, "message": limitErr.Error()})
	}

	events, err := e.usecases.RetrieveEvents(spectrumId, limit, page)

	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Request was successful", "events": events})
	}
}

func (e EventHandlers) Register(ctx *gin.Context) {
	var entity entity.EventRegistrationFromJsonEntity

	canGo := pkg.ParseJson(ctx, &entity)
	if !canGo {
		return
	}

	details, err := e.usecases.Register(entity)
	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "User has been registered", "event_registration": details})
}

func (e EventHandlers) ChangeRegistrationStatus(ctx *gin.Context) {
	var entity entity.ChangeRegistrationStatusEntity

	canGo := pkg.ParseJson(ctx, &entity)
	if !canGo {
		return
	}

	newRegistration, err := e.usecases.ChangeRegistrationStatus(entity)
	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Status has been updated", "event_registration": newRegistration})
}

func (e EventHandlers) WithdrawRegistration(ctx *gin.Context) {

	registrationId, exists := ctx.GetQuery("registration_id")
	if !exists {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": "Provide registration id"})
		return
	}

	err := e.usecases.WithdrawRegistration(registrationId)

	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Registration has been withdrawn"})
}
