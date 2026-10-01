package routes

import (
	"espectro/enums"
	"espectro/handlers"
	"espectro/middlewares"
	repositoryimple "espectro/repository_imple"
	"espectro/usecases"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterCompetitionRoutes(r *gin.RouterGroup, db *gorm.DB, cld *cloudinary.Cloudinary) {

	competitionRepo := repositoryimple.NewCompetitionPostgresRepo(db)
	mediaRepo := repositoryimple.NewMediaCloudinaryRepo(cld)
	venueRepo := repositoryimple.NewVenuePostgresRepo(db)
	prizeRepo := repositoryimple.NewPrizePostgresRepo(db)
	adminRepo := repositoryimple.NewAdminPostgresRepo(db)
	eventRepo := repositoryimple.NewEventPostgresRepo(db)
	userRepo := repositoryimple.NewUserPostgresRepo(db)

	adminUsecases := usecases.NewAdminUsecases(adminRepo)
	userUsecases := usecases.NewUserUsecases(userRepo)
	competitionUsecases := usecases.NewCompetitionUsecases(
		competitionRepo,
		venueRepo,
		prizeRepo,
		eventRepo,
		mediaRepo,
	)

	leaderMemberAdminMiddleware := middlewares.NewAdminMiddleWare(adminUsecases, enums.LeaderAndMemberMiddleware)
	userAdminMiddleware := middlewares.NewUserAdminMiddleware(userUsecases, adminUsecases)

	competitionHandlers := handlers.NewCompetitionHandlers(competitionUsecases)

	competitionApi := r.Group("/competition")

	competitionApi.POST("", leaderMemberAdminMiddleware.AdminMiddleWare, competitionHandlers.CreateCompetition)
	competitionApi.PATCH("", leaderMemberAdminMiddleware.AdminMiddleWare, competitionHandlers.UpdateCompetition)
	competitionApi.DELETE("", leaderMemberAdminMiddleware.AdminMiddleWare, competitionHandlers.DeleteCompetition)
	competitionApi.GET("", userAdminMiddleware.UserAdminMiddleware, competitionHandlers.RetrieveCompetitions)
}
