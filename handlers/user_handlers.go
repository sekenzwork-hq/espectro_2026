package handlers

import (
	"espectro/usecases"

	"github.com/gin-gonic/gin"
)

type UserHandlers struct {
	usecases usecases.UserUsecases
}

func NewUserHandlers(usecases usecases.UserUsecases) UserHandlers {
	return UserHandlers{
		usecases: usecases,
	}
}

func (h *UserHandlers) RegisterUser(ctx *gin.Context) {

	ctx.JSON(200, gin.H{
		"message": "connected",
	})
}
