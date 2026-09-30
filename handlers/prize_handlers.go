package handlers

import (
	"espectro/entity"
	"espectro/enums"
	"espectro/pkg"
	"espectro/usecases"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PrizeHandlers struct {
	prizeUsecases usecases.PrizeUsecases
}

func NewPrizeHandlers(prizeUsecase usecases.PrizeUsecases) PrizeHandlers {
	return PrizeHandlers{prizeUsecases: prizeUsecase}
}

func (p PrizeHandlers) CreatePrize(ctx *gin.Context) {

	form, formErr := ctx.MultipartForm()

	if formErr != nil {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": "Invalid form data"})
		return
	}

	title, titleExists := ctx.GetPostForm("title")
	des, desExists := ctx.GetPostForm("description")
	eventId, eventIdExists := ctx.GetPostForm("event_id")
	prizeType, prizeTypeExists := ctx.GetPostForm("type")
	amountStr, amountExists := ctx.GetPostForm("amount")

	var amount float32 = 0

	var logo *multipart.FileHeader
	if len(form.File["logo"]) >= 1 {
		fmt.Println("Exists")
		logo = form.File["logo"][0]
	}

	if !titleExists {
		ctx.JSON(http.StatusNotAcceptable, pkg.EmptyMessage("Provide title"))
	} else if !desExists {
		ctx.JSON(http.StatusNotAcceptable, pkg.EmptyMessage("Provide description"))
	} else if !eventIdExists {
		ctx.JSON(http.StatusNotAcceptable, pkg.EmptyMessage("Provide event id"))
	} else if !prizeTypeExists {
		ctx.JSON(http.StatusNotAcceptable, pkg.EmptyMessage("Provide type"))
	} else if amountExists {

		floatAmount, parseErr := strconv.ParseFloat(amountStr, 32)
		if parseErr != nil {
			ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": "Invalid amount"})
		}
		amount = float32(floatAmount)
	} else {

		createdPrize, creationErr := p.prizeUsecases.CreatePrize(entity.PrizeCreateEntity{
			Title:       title,
			Description: des, Amount: amount,
			EventId: eventId,
			Logo:    logo,
			Type:    enums.PrizeTypeEnum(prizeType),
		})

		if creationErr != nil {
			code := pkg.GetStatusCodeForError(creationErr)
			ctx.JSON(code, gin.H{"status": code, "message": creationErr.Error()})
		} else {
			ctx.JSON(http.StatusCreated, gin.H{"status": 201, "message": "Prize has been created successfully", "prize": createdPrize})
		}
	}
}

func (p PrizeHandlers) UpdatePrize(ctx *gin.Context) {

	form, formErr := ctx.MultipartForm()

	if formErr != nil {
		ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": "Invalid form data"})
		return
	}

	id, idExists := ctx.GetPostForm("prize_id")

	if !idExists {
		ctx.JSON(http.StatusNotAcceptable, pkg.EmptyMessage("Provide prize id"))
	}

	titleForm, titleExists := ctx.GetPostForm("title")
	desForm, desExists := ctx.GetPostForm("description")
	eventIdForm, eventIdExists := ctx.GetPostForm("event_id")
	prizeTypeForm, prizeTypeExists := ctx.GetPostForm("type")
	amountStrForm, amountExists := ctx.GetPostForm("amount")

	var title *string
	var des *string
	var eventId *string
	var prizeType *enums.PrizeTypeEnum
	var amount *float32
	var logo *multipart.FileHeader

	if titleExists {
		title = &titleForm
	}
	if desExists {
		des = &desForm
	}
	if eventIdExists {
		eventId = &eventIdForm
	}
	if prizeTypeExists {
		prizeType = (*enums.PrizeTypeEnum)(&prizeTypeForm)
	}

	if amountExists {
		floatAmount, parseErr := strconv.ParseFloat(amountStrForm, 32)
		if parseErr != nil {
			ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": "Invalid amount"})
		}
		amountPtr := (float32)(floatAmount)
		amount = &amountPtr
	}

	if len(form.File["logo"]) >= 1 {
		logo = form.File["logo"][0]
	}

	updatedPrize, updationErr := p.prizeUsecases.UpdatePrize(id, entity.PrizeUpdateEntity{
		Title:       title,
		Description: des,
		Amount:      amount,
		EventId:     eventId,
		Logo:        logo,
		Type:        prizeType,
	})

	if updationErr != nil {
		code := pkg.GetStatusCodeForError(updationErr)
		ctx.JSON(code, gin.H{"status": code, "message": updationErr.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Prize has been updated successfully", "prize": updatedPrize})
	}

}

func (p PrizeHandlers) DeletePrize(ctx *gin.Context) {

	prizeId, exists := ctx.GetQuery("prize_id")

	if !exists {
		ctx.JSON(http.StatusNotAcceptable, pkg.EmptyMessage("Provide prize id"))
		return
	}

	err := p.prizeUsecases.DeletePrize(prizeId)

	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
	} else {
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Prize has been deleted successfully"})
	}
}
