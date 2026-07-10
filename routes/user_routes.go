package routes

import (
	"espectro/handlers"
	repositoryimple "espectro/repository_imple"
	"espectro/usecases"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterUserRoutes(r *gin.RouterGroup, db *gorm.DB) {

	usersApi := r.Group("/user")

	repo := repositoryimple.NewUserPostgresRepo(db)

	usecases := usecases.NewUserUsecases(repo)

	userHandlers := handlers.NewUserHandlers(usecases)

	usersApi.POST("/register", userHandlers.RegisterUser)

}
