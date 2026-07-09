package middlewares

import (
	"espectro/pkg"
	"espectro/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdminMiddleWare struct {
	usecases usecases.AdminUsecases
}

func NewAdminMiddleWare(usecases usecases.AdminUsecases) AdminMiddleWare {
	return AdminMiddleWare{usecases: usecases}
}

func (a AdminMiddleWare) AdminMiddleWare(ctx *gin.Context) {

	headerToken := ctx.GetHeader("Authorization")

	if len(headerToken) == 0 {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": 401, "message": "Invalid token"})
		return
	}

	parsedAdminId, parseErr := pkg.ParseJWTFromAdmin(headerToken)

	if parseErr != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": 401, "message": "Invalid token"})
		return
	}

	err := a.usecases.CheckAdminExists(parsedAdminId)

	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.AbortWithStatusJSON(code, gin.H{"status": code, "message": err.Error()})
		return
	}
	ctx.Set("admin_id", parsedAdminId)
	ctx.Next()

}
