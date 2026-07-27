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
	adminUsecases := usecases.NewAdminUsecases(adminRepo)
	orgUsecases := usecases.NewOrganizationUsecases(orgRepo, mediaRepo)
	orgHandlers := handlers.NewOrganizationHandlers(orgUsecases)

	leaderAndMemberMiddleware := middlewares.NewAdminMiddleWare(adminUsecases, enums.LeaderAndMemberMiddleware)

	orgApi := r.Group("/organization")

	orgApi.POST("/create", orgHandlers.CreateOrganization)
	orgApi.PATCH("/update", orgHandlers.UpdateOrganization)
	orgApi.DELETE("/delete", orgHandlers.DeleteOrganization)
	orgApi.GET("", leaderAndMemberMiddleware.AdminMiddleWare, orgHandlers.RetrieveOrganizations)
}
