package middlewares

import (
	customerrors "espectro/custom_errors"
	"espectro/enums"
	"espectro/pkg"
	"espectro/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdminMiddleWare struct {
	adminType enums.AdminMiddlewareType
	usecases  usecases.AdminUsecases
}

func NewAdminMiddleWare(usecases usecases.AdminUsecases, adminType enums.AdminMiddlewareType) AdminMiddleWare {
	return AdminMiddleWare{usecases: usecases, adminType: adminType}
}

func (a AdminMiddleWare) AdminMiddleWare(ctx *gin.Context) {

	if !a.adminType.IsValid() {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": 500, "message": "Something went wrong while operating"})
	}

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

	var err error

	switch a.adminType {
	case enums.AllAdminMiddleware:

		err = a.usecases.CheckAdminExists(parsedAdminId)

	case enums.LeaderMiddleware:

		role, roleErr := a.usecases.RetrieveAdminRoleByID(parsedAdminId)

		if roleErr != nil {
			err = roleErr
		} else if role != enums.Leader {
			err = &customerrors.PermissionError{OrgError: "Current admin doesn't have permission"}
		}
	}

	if err != nil {
		code := pkg.GetStatusCodeForError(err)
		ctx.AbortWithStatusJSON(code, gin.H{"status": code, "message": err.Error()})
		return
	}
	ctx.Set("admin_id", parsedAdminId)
	ctx.Next()

}
