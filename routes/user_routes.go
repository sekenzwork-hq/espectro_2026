package routes

import (
	"espectro/handlers"
	"espectro/repository"
	"espectro/usecases"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterUserRoutes(r *gin.Engine, db *gorm.DB) {

	api := r.Group("/api/v1")

	usersApi := api.Group("/user")

	repo := repository.NewUserPostgresRepo(db)

	usecases := usecases.NewUserUsecases(repo)

	userHandlers := handlers.NewUserHandlers(usecases)

	usersApi.POST("/register", userHandlers.RegisterUser)

}
