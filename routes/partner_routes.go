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

func RegisterPartnerRoutes(r *gin.RouterGroup, db *gorm.DB) {

	partnerRepo := repositoryimple.NewPartnerPostgres(db)
	partnerUsecases := usecases.NewPartnerUsecases(partnerRepo)
	adminRepo := repositoryimple.NewAdminPostgresRepo(db)
	adminUsecases := usecases.NewAdminUsecases(adminRepo)
	handlers := handlers.NewPartnerHandlers(partnerUsecases)

	leaderAndMemberMiddleware := middlewares.NewAdminMiddleWare(adminUsecases, enums.LeaderAndMemberMiddleware)

	partnerApi := r.Group("/partner")

	partnerApi.POST("/create", leaderAndMemberMiddleware.AdminMiddleWare, handlers.AddPartner)
	partnerApi.PATCH("/update", leaderAndMemberMiddleware.AdminMiddleWare, handlers.UpdatePartner)
	partnerApi.DELETE("/delete", leaderAndMemberMiddleware.AdminMiddleWare, handlers.DeletePartner)
}
