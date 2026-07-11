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
	adminApi.POST("/create", leaderAdminProtectedMiddleware.AdminMiddleWare, handlers.CreateNewAdmin)
	//This route is for deleting another admin whose role is member or volunteer
	adminApi.DELETE("/delete", leaderAdminProtectedMiddleware.AdminMiddleWare, handlers.DeleteMemberOrVolunteer)
	//This route is for updating current admin's details
	adminApi.PATCH("/update", allAdminProtectedMiddleware.AdminMiddleWare, handlers.UpdateCurrentAdmin)
	//This route is for updating another admin's role
	adminApi.PATCH("/update/role", leaderAdminProtectedMiddleware.AdminMiddleWare, handlers.UpdateAdminRole)
}
