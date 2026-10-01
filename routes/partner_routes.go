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

func RegisterPartnerRoutes(r *gin.RouterGroup, db *gorm.DB, cld *cloudinary.Cloudinary) {

	partnerRepo := repositoryimple.NewPartnerPostgres(db)
	mediaRepo := repositoryimple.NewMediaCloudinaryRepo(cld)
	adminRepo := repositoryimple.NewAdminPostgresRepo(db)
	userRepo := repositoryimple.NewUserPostgresRepo(db)

	adminUsecases := usecases.NewAdminUsecases(adminRepo)
	partnerUsecases := usecases.NewPartnerUsecases(partnerRepo, mediaRepo)
	userUsecases := usecases.NewUserUsecases(userRepo)

	handlers := handlers.NewPartnerHandlers(partnerUsecases)

	leaderAndMemberMiddleware := middlewares.NewAdminMiddleWare(adminUsecases, enums.LeaderAndMemberMiddleware)
	userAdminMiddleware := middlewares.NewUserAdminMiddleware(userUsecases, adminUsecases)

	partnerApi := r.Group("/partner")

	partnerApi.POST("", leaderAndMemberMiddleware.AdminMiddleWare, handlers.CreatePartner)
	partnerApi.PATCH("", leaderAndMemberMiddleware.AdminMiddleWare, handlers.UpdatePartner)
	partnerApi.DELETE("", leaderAndMemberMiddleware.AdminMiddleWare, handlers.DeletePartner)
	partnerApi.GET("", userAdminMiddleware.UserAdminMiddleware, handlers.RetrievePartner)
}
