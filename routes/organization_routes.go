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

func RegisterOrganizationRoutes(r *gin.RouterGroup, db *gorm.DB, cld *cloudinary.Cloudinary) {

	orgRepo := repositoryimple.NewOrganizationPostgresRepo(db)
	mediaRepo := repositoryimple.NewMediaCloudinaryRepo(cld)
	adminRepo := repositoryimple.NewAdminPostgresRepo(db)
	userRepo := repositoryimple.NewUserPostgresRepo(db)

	orgUsecases := usecases.NewOrganizationUsecases(orgRepo, mediaRepo)
	adminUsecases := usecases.NewAdminUsecases(adminRepo)
	userUsecases := usecases.NewUserUsecases(userRepo)

	orgHandlers := handlers.NewOrganizationHandlers(orgUsecases)

	leaderAndMemberMiddleware := middlewares.NewAdminMiddleWare(adminUsecases, enums.LeaderAndMemberMiddleware)
	userMiddleware := middlewares.NewUserMiddleware(userUsecases)

	orgApi := r.Group("/organization")

	orgApi.POST("", userMiddleware.UserMiddleware, orgHandlers.CreateOrganization)
	orgApi.PATCH("", userMiddleware.UserMiddleware, orgHandlers.UpdateOrganization)
	orgApi.DELETE("", userMiddleware.UserMiddleware, orgHandlers.DeleteOrganization)
	orgApi.GET("", leaderAndMemberMiddleware.AdminMiddleWare, orgHandlers.RetrieveOrganizations)
}
