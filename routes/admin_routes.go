package routes

import (
	"espectro/handlers"
	"espectro/middlewares"
	"espectro/repository"
	"espectro/usecases"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterAdminRoutes(r *gin.RouterGroup, db *gorm.DB) {

	repo := repository.NewAdminPostgresRepo(db)

	usecases := usecases.NewAdminUsecases(repo)

	handlers := handlers.NewAdminHandlers(usecases)

	adminMiddleware := middlewares.NewAdminMiddleWare(usecases)

	adminApi := r.Group("/admin")

	adminApi.POST("/login", handlers.Login)
	adminApi.POST("/create", adminMiddleware.AdminMiddleWare, handlers.CreateNewAdmin)
	adminApi.DELETE("/delete", adminMiddleware.AdminMiddleWare, handlers.DeleteMemberOrVolunteer)
	adminApi.PATCH("/update", adminMiddleware.AdminMiddleWare, handlers.UpdateCurrentAdmin)
}
