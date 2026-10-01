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
	adminRepo := repositoryimple.NewAdminPostgresRepo(db)
	userRepo := repositoryimple.NewUserPostgresRepo(db)

	adminUsecases := usecases.NewAdminUsecases(adminRepo)
	userUsecases := usecases.NewUserUsecases(userRepo)

	leaderAndMemberMiddleware := middlewares.NewAdminMiddleWare(adminUsecases, enums.LeaderAndMemberMiddleware)
	userAdminMiddleware := middlewares.NewUserAdminMiddleware(userUsecases, adminUsecases)

	handlers := handlers.NewVenueHandlers(venueUsecases)

	venueApi := r.Group("/venue")

	venueApi.POST("", leaderAndMemberMiddleware.AdminMiddleWare, handlers.CreateVenue)
	venueApi.DELETE("", leaderAndMemberMiddleware.AdminMiddleWare, handlers.DeleteVenue)
	venueApi.PATCH("", leaderAndMemberMiddleware.AdminMiddleWare, handlers.UpdateVenue)
	venueApi.GET("", userAdminMiddleware.UserAdminMiddleware, handlers.RetrieveVenue)
}
