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

func RegisterInvestorRoutes(r *gin.RouterGroup, db *gorm.DB, cld *cloudinary.Cloudinary) {

	investorRepo := repositoryimple.NewInvestorPostgresRepo(db)
	adminRepo := repositoryimple.NewAdminPostgresRepo(db)
	mediaRepo := repositoryimple.NewMediaCloudinaryRepo(cld)
	investorUsecases := usecases.NewInvestorUsecases(investorRepo, mediaRepo)
	adminUsecases := usecases.NewAdminUsecases(adminRepo)
	leaderAndMemberMiddleware := middlewares.NewAdminMiddleWare(adminUsecases, enums.LeaderAndMemberMiddleware)
	handlers := handlers.NewInvestorHandlers(investorUsecases)

	investorApi := r.Group("investor")

	investorApi.POST("/add", leaderAndMemberMiddleware.AdminMiddleWare, handlers.AddInvestor)
	investorApi.PATCH("/update", leaderAndMemberMiddleware.AdminMiddleWare, handlers.UpdateInvestor)
}
