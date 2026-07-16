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
	partnerUsecases := usecases.NewPartnerUsecases(partnerRepo, mediaRepo)
	adminRepo := repositoryimple.NewAdminPostgresRepo(db)
	adminUsecases := usecases.NewAdminUsecases(adminRepo)
	handlers := handlers.NewPartnerHandlers(partnerUsecases)

	leaderAndMemberMiddleware := middlewares.NewAdminMiddleWare(adminUsecases, enums.LeaderAndMemberMiddleware)

	partnerApi := r.Group("/partner")

	partnerApi.POST("/add", leaderAndMemberMiddleware.AdminMiddleWare, handlers.AddPartner)
	partnerApi.PATCH("/update", leaderAndMemberMiddleware.AdminMiddleWare, handlers.UpdatePartner)
	partnerApi.DELETE("/delete", leaderAndMemberMiddleware.AdminMiddleWare, handlers.DeletePartner)
	partnerApi.GET("", handlers.RetrievePartner)
}
