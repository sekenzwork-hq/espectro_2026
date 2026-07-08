package handlers

import (
	"errors"
	customerrors "espectro/custom_errors"
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

	id, validationOrDBError := h.usecases.RegisterUser(userEntity)

	if validationOrDBError != nil {

		var validationError *customerrors.ValidationError
		isValidationError := errors.As(validationOrDBError, &validationError)

		if isValidationError {
			ctx.JSON(http.StatusNotAcceptable, gin.H{"status": 406, "message": validationOrDBError.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": 500, "message": validationOrDBError.Error()})
		}
	} else {
		ctx.JSON(http.StatusOK, gin.H{"status": 200, "message": "User registered", "id": id})
	}
}
