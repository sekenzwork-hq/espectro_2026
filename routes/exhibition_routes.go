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

func RegisterExhibtionRoutes(r *gin.RouterGroup, db *gorm.DB, cld *cloudinary.Cloudinary) {

	exhibtionRepo := repositoryimple.NewExhibitionPostgresRepo(db)
	eventRepo := repositoryimple.NewEventPostgresRepo(db)
	mediaRepo := repositoryimple.NewMediaCloudinaryRepo(cld)
	organizationRepo := repositoryimple.NewOrganizationPostgresRepo(db)
	adminRepo := repositoryimple.NewAdminPostgresRepo(db)
	staffRepo := repositoryimple.NewStaffPostgresRepo(db)

	exhibitionUsecases := usecases.NewExhibitionUsecases(exhibtionRepo, mediaRepo, eventRepo, organizationRepo, adminRepo, staffRepo)
	adminUsecases := usecases.NewAdminUsecases(adminRepo)

	exhibtionHandlers := handlers.NewExhibitionHandlers(exhibitionUsecases)

	memberLeaderAdminMiddleware := middlewares.NewAdminMiddleWare(adminUsecases, enums.LeaderAndMemberMiddleware)

	exhibitionApi := r.Group("/exhibition")

	exhibitionApi.POST("", exhibtionHandlers.CreateExhibition)
	exhibitionApi.PATCH("/user", exhibtionHandlers.UpdateExhibitionFromUserSide)
	exhibitionApi.PATCH("/admin", memberLeaderAdminMiddleware.AdminMiddleWare, exhibtionHandlers.UpdateExhibitionFromAdminSide)
	exhibitionApi.DELETE("/admin", memberLeaderAdminMiddleware.AdminMiddleWare, exhibtionHandlers.DeleteExhibition)
	exhibitionApi.DELETE("/user", exhibtionHandlers.DeleteExhibition)

}
