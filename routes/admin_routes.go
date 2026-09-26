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

func RegisterAdminRoutes(r *gin.RouterGroup, db *gorm.DB) {

	repo := repositoryimple.NewAdminPostgresRepo(db)

	usecases := usecases.NewAdminUsecases(repo)

	handlers := handlers.NewAdminHandlers(usecases)

	allAdminProtectedMiddleware := middlewares.NewAdminMiddleWare(usecases, enums.AllAdminMiddleware)
	leaderAdminProtectedMiddleware := middlewares.NewAdminMiddleWare(usecases, enums.LeaderMiddleware)

	adminApi := r.Group("/admin")

	adminApi.POST("/login", handlers.Login)
	adminApi.POST("", leaderAdminProtectedMiddleware.AdminMiddleWare, handlers.CreateNewAdmin)
	adminApi.DELETE("", leaderAdminProtectedMiddleware.AdminMiddleWare, handlers.DeleteMemberOrVolunteer)
	adminApi.PATCH("", allAdminProtectedMiddleware.AdminMiddleWare, handlers.UpdateCurrentAdmin)
	adminApi.PATCH("/role", leaderAdminProtectedMiddleware.AdminMiddleWare, handlers.UpdateAdminRole)
}
