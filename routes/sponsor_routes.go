package routes

import (
	"espectro/database"
	"espectro/enums"
	"espectro/handlers"
	"espectro/middlewares"
	repositoryimple "espectro/repository_imple"
	"espectro/usecases"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterSponsorRoutes(r *gin.RouterGroup, db *gorm.DB, cld *cloudinary.Cloudinary) {

	sponsorRepo := repositoryimple.NewSponsorPostgresRepo(db)
	eventRepo := repositoryimple.NewEventPostgresRepo(db)
	eventSponsorRepo := repositoryimple.NewEventSponsorPostgresRepo(db)
	adminRepo := repositoryimple.NewAdminPostgresRepo(db)
	mediaRepo := repositoryimple.NewMediaCloudinaryRepo(cld)
	userRepo := repositoryimple.NewUserPostgresRepo(db)

	transactionManager := database.NewTransactionManager(db)

	sponsorUsecases := usecases.NewSponsorUsecases(sponsorRepo, eventSponsorRepo, eventRepo, transactionManager, mediaRepo)
	adminUsecases := usecases.NewAdminUsecases(adminRepo)
	userUsecases := usecases.NewUserUsecases(userRepo)

	leaderAndMemberMiddleware := middlewares.NewAdminMiddleWare(adminUsecases, enums.LeaderAndMemberMiddleware)
	userAdminMiddleware := middlewares.NewUserAdminMiddleware(userUsecases, adminUsecases)

	handlers := handlers.NewSponsorHandlers(sponsorUsecases)

	sponsorApi := r.Group("sponsor")

	sponsorApi.POST("", leaderAndMemberMiddleware.AdminMiddleWare, handlers.CreateSponsor)
	sponsorApi.PATCH("", leaderAndMemberMiddleware.AdminMiddleWare, handlers.UpdateSponsor)
	sponsorApi.DELETE("", leaderAndMemberMiddleware.AdminMiddleWare, handlers.DeleteSponsor)
	sponsorApi.GET("", userAdminMiddleware.UserAdminMiddleware, handlers.RetrieveSponsors)
}
