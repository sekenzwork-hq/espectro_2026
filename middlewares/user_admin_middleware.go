package middlewares

import (
	"espectro/pkg"
	"espectro/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserAdminMiddleware struct {
	userUsecases  usecases.UserUsecases
	adminUsecases usecases.AdminUsecases
}

func NewUserAdminMiddleware(userUsecases usecases.UserUsecases, adminUsecases usecases.AdminUsecases) UserAdminMiddleware {
	return UserAdminMiddleware{userUsecases: userUsecases, adminUsecases: adminUsecases}
}

func (u UserAdminMiddleware) UserAdminMiddleware(ctx *gin.Context) {

	headerToken := ctx.GetHeader("Authorization")

	if len(headerToken) == 0 {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": 401, "message": "Invalid token"})
		return
	}

	parsedAdminId, parseAdminIdErr := pkg.ParseJWTFromAdmin(headerToken)
	parsedUserId, parseUserIdErr := pkg.ParseJWTFromUser(headerToken)

	if pkg.ValidateUUID(parsedAdminId) {
		if parseAdminIdErr != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": 401, "message": "Invalid token"})
			return
		}

		err := u.adminUsecases.AdminExists(parsedAdminId)

		if err != nil {
			code := pkg.GetStatusCodeForError(err)
			ctx.AbortWithStatusJSON(code, gin.H{"status": code, "message": err.Error()})
		} else {
			ctx.Set("admin_id", parsedAdminId)
			ctx.Next()
		}

	} else if pkg.ValidateUUID(parsedUserId) {
		if parseUserIdErr != nil {
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
	} else {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": 401, "message": "Invalid token"})
	}
}
