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

func RegisterPrizeRoutes(r *gin.RouterGroup, db *gorm.DB, cld *cloudinary.Cloudinary) {

	prizeRepo := repositoryimple.NewPrizePostgresRepo(db)
	eventRepo := repositoryimple.NewEventPostgresRepo(db)
	adminRepo := repositoryimple.NewAdminPostgresRepo(db)
	mediaRepo := repositoryimple.NewMediaCloudinaryRepo(cld)

	prizeUsecases := usecases.NewPrizeUsecases(prizeRepo, mediaRepo, eventRepo)
	adminUsecases := usecases.NewAdminUsecases(adminRepo)

	prizeHandlers := handlers.NewPrizeHandlers(prizeUsecases)

	leaderMemberAdminMiddleware := middlewares.NewAdminMiddleWare(adminUsecases, enums.LeaderAndMemberMiddleware)

	prizeApi := r.Group("/prize")

	prizeApi.POST("", leaderMemberAdminMiddleware.AdminMiddleWare, prizeHandlers.CreatePrize)
	prizeApi.PATCH("", leaderMemberAdminMiddleware.AdminMiddleWare, prizeHandlers.UpdatePrize)
	prizeApi.DELETE("", leaderMemberAdminMiddleware.AdminMiddleWare, prizeHandlers.DeletePrize)

}
