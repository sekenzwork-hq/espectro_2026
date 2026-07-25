package routes

import (
	"espectro/enums"
	"espectro/handlers"
	"espectro/middlewares"
	repositoryimple "espectro/repository_imple"
	"espectro/usecases"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterStaffRoutes(r *gin.RouterGroup, db *gorm.DB) {

	staffRepo := repositoryimple.NewStaffPostgresRepo(db)
	adminRepo := repositoryimple.NewAdminPostgresRepo(db)
	adminUsecases := usecases.NewAdminUsecases(adminRepo)
	staffUsecases := usecases.NewStaffUsecases(staffRepo)
	handlers := handlers.NewStaffHandlers(staffUsecases)

	adminMiddleware := middlewares.NewAdminMiddleWare(adminUsecases, enums.LeaderAndMemberMiddleware)

	staffApi := r.Group("/staff")
	staffApi.POST("/create", adminMiddleware.AdminMiddleWare, handlers.CreateStaff)
	staffApi.PATCH("/update", adminMiddleware.AdminMiddleWare, handlers.UpdateStaff)
}
