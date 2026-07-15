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

	var entity entity.EventFromJsonEntity

	canGo := pkg.ParseJson(ctx, &entity)

	if !canGo {
		return
	}

	event, err := e.usecases.CreateEvent(entity)

	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Event has been created", "id": event.Id, "created_at": event.CreatedAt})
	}
}

func (e EventHandlers) UpdateEvent(ctx *gin.Context) {
	var entity entity.EventUpdateEntity

	canGo := pkg.ParseJson(ctx, &entity)

	if !canGo {
		return
	}

	err := e.usecases.UpdateEvent(entity)

	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Event has been updated"})
	}

}

func (e EventHandlers) DeleteEvent(ctx *gin.Context) {

	eventId := ctx.Query("event_id")

	err := e.usecases.DeleteEvent(eventId)

	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Event has been deleted"})
	}
}

func (e EventHandlers) RetrieveEvents(ctx *gin.Context) {

	limit, page := pkg.ParsePageAndLimit(ctx.Query("limit"), ctx.Query("page"))

	events, err := e.usecases.RetrieveEvents(limit, page)

	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Request was successful", "events": events})
	}
}
