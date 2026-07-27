package routes

import (
	"espectro/handlers"
	repositoryimple "espectro/repository_imple"
	"espectro/usecases"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterOrganizationRoutes(r *gin.RouterGroup, db *gorm.DB, cld *cloudinary.Cloudinary) {

	orgRepo := repositoryimple.NewOrganizationPostgresRepo(db)
	mediaRepo := repositoryimple.NewMediaCloudinaryRepo(cld)
	orgUsecases := usecases.NewOrganizationUsecases(orgRepo, mediaRepo)
	orgHandlers := handlers.NewOrganizationHandlers(orgUsecases)

	orgApi := r.Group("/organization")

	orgApi.POST("/create", orgHandlers.CreateOrganization)
}
