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

	adminUsecases := usecases.NewAdminUsecases(adminRepo)
	competitionUsecases := usecases.NewCompetitionUsecases(competitionRepo, venueRepo, prizeRepo, mediaRepo)

	leaderMemberAdminMiddleware := middlewares.NewAdminMiddleWare(adminUsecases, enums.LeaderAndMemberMiddleware)

	competitionHandlers := handlers.NewCompetitionHandlers(competitionUsecases)

	competitionApi := r.Group("/competition")

	competitionApi.POST("", leaderMemberAdminMiddleware.AdminMiddleWare, competitionHandlers.CreateCompetition)
}
