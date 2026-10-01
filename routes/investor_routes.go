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
	userRepo := repositoryimple.NewUserPostgresRepo(db)

	investorUsecases := usecases.NewInvestorUsecases(investorRepo, mediaRepo)
	adminUsecases := usecases.NewAdminUsecases(adminRepo)
	userUsecases := usecases.NewUserUsecases(userRepo)

	leaderAndMemberMiddleware := middlewares.NewAdminMiddleWare(adminUsecases, enums.LeaderAndMemberMiddleware)
	userAdminMiddleware := middlewares.NewUserAdminMiddleware(userUsecases, adminUsecases)

	handlers := handlers.NewInvestorHandlers(investorUsecases)

	investorApi := r.Group("investor")

	investorApi.POST("", leaderAndMemberMiddleware.AdminMiddleWare, handlers.CreateInvestor)
	investorApi.PATCH("", leaderAndMemberMiddleware.AdminMiddleWare, handlers.UpdateInvestor)
	investorApi.DELETE("", leaderAndMemberMiddleware.AdminMiddleWare, handlers.DeleteInvestor)
	investorApi.GET("", userAdminMiddleware.UserAdminMiddleware, handlers.RetrieveInvestors)
}
