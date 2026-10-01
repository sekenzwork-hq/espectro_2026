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

func (u UserHandlers) RegisterUser(ctx *gin.Context) {

	var userEntity entity.UserJsonCreateEntity

	canGo := pkg.ParseJson(ctx, &userEntity)

	if !canGo {
		return
	}

	user, err := u.usecases.RegisterUser(userEntity)

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

func (u UserHandlers) Login(ctx *gin.Context) {

	var entity entity.UserCredentialsJsonEntity

	canGo := pkg.ParseJson(ctx, &entity)
	if !canGo {
		return
	}

	userId, err := u.usecases.Login(entity)

	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.JSON(code, gin.H{"status": code, "message": err.Error()})
	} else {

		token, genErr := pkg.GenerateJWTForUser(userId)

		if genErr != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": 500, "message": "Something went wrong while operating"})
		} else {
			ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "Logged successfully", "access_token": token})
		}
	}
}
