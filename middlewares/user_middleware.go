package middlewares

import (
	"espectro/pkg"
	"espectro/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserMiddleware struct {
	userUsecases usecases.UserUsecases
}

func NewUserMiddleware(userUsecases usecases.UserUsecases) UserMiddleware {
	return UserMiddleware{userUsecases: userUsecases}
}

func (u UserMiddleware) UserMiddleware(ctx *gin.Context) {

	headerToken := ctx.GetHeader("Authorization")

	if len(headerToken) == 0 {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": 401, "message": "Invalid token"})
		return
	}

	parsedUserId, parseErr := pkg.ParseJWTFromUser(headerToken)

	if parseErr != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": 401, "message": "Invalid token"})
		return
	}

	err := u.userUsecases.CheckUserExists(parsedUserId)

	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.AbortWithStatusJSON(code, gin.H{"status": code, "message": err.Error()})
	} else {
		ctx.Set("user_id", parsedUserId)
		ctx.Next()
	}

}
