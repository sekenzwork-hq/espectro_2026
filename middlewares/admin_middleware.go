package middlewares

import (
	"espectro/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminMiddleWare(ctx *gin.Context) {

	headerToken := ctx.GetHeader("Authorization")

	if len(headerToken) == 0 {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": 401, "message": "Invalid token"})
		return
	}

	parsedAdminId, parseErr := pkg.ParseJWTFromAdmin(headerToken)

	if parseErr != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": 401, "message": "Invalid token"})
	} else {
		ctx.Set("admin_id", parsedAdminId)
		ctx.Next()
	}
}
