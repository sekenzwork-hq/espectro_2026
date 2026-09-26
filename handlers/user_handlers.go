package handlers

import (
	"espectro/entity"
	"espectro/pkg"
	"espectro/usecases"
	"net/http"

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

	var userEntity entity.UserEntity

	canGo := pkg.ParseJson(ctx, &userEntity)

	if !canGo {
		return
	}

	user, err := h.usecases.RegisterUser(userEntity)

	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
	} else {
		token, tokenErr := pkg.GenerateJWTForUser(user.Id)
		if tokenErr != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": 500, "message": "Something went wrong while operating"})
		} else {
			ctx.JSON(http.StatusCreated, gin.H{"status": 201, "message": "User has been registered", "user": user, "access_token": token})
		}

	}
}
