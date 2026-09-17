package routes

import (
	"espectro/handlers"
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

	exhibitionUsecases := usecases.NewExhibitionUsecases(exhibtionRepo, mediaRepo, eventRepo, organizationRepo)

	exhibtionHandlers := handlers.NewExhibitionHandlers(exhibitionUsecases)

	exhibitionApi := r.Group("/exhibition")

	exhibitionApi.POST("", exhibtionHandlers.CreateExhibition)

}
