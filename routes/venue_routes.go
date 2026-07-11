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

func RegisterVenueRoutes(r *gin.RouterGroup, db *gorm.DB) {

	venueRepo := repositoryimple.NewVenuePostgresRepo(db)
	venueUsecases := usecases.NewVenueUsecases(venueRepo)
	handlers := handlers.NewVenueHandlers(venueUsecases)
	adminRepo := repositoryimple.NewAdminPostgresRepo(db)
	adminUsecase := usecases.NewAdminUsecases(adminRepo)

	leaderAndMemberMiddleware := middlewares.NewAdminMiddleWare(adminUsecase, enums.LeaderAndMemberMiddleware)

	venueApi := r.Group("/venue")

	venueApi.POST("/create", leaderAndMemberMiddleware.AdminMiddleWare, handlers.CreateVenue)
}
