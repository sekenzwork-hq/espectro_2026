package routes

import (
	"espectro/handlers"
	"espectro/repository"
	"espectro/usecases"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterUserRoutes(r *gin.RouterGroup, db *gorm.DB) {

	usersApi := r.Group("/user")

	repo := repository.NewUserPostgresRepo(db)

	usecases := usecases.NewUserUsecases(repo)

	userHandlers := handlers.NewUserHandlers(usecases)

	usersApi.POST("/register", userHandlers.RegisterUser)

}
