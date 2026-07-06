package routes

import (
	"espectro/internal/delivery/http/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine) {

	api := r.Group("/api/v1")

	usersApi := api.Group("/user")

	usersApi.POST("/register", handlers.RegisterUser)

}
