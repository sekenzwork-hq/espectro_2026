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

func RegisterSpectrumRoutes(r *gin.RouterGroup, db *gorm.DB, cld *cloudinary.Cloudinary) {

	spectrumRepo := repositoryimple.NewSpectrumPostgresRepo(db)
	adminRepo := repositoryimple.NewAdminPostgresRepo(db)
	cldMediaRepo := repositoryimple.NewMediaCloudinaryRepo(cld)
	spectrumUsecase := usecases.NewSpectrumUsecases(spectrumRepo, cldMediaRepo)
	adminUsecase := usecases.NewAdminUsecases(adminRepo)
	handlers := handlers.NewSpectrumHandlers(spectrumUsecase)

	leaderAndMemberMiddleware := middlewares.NewAdminMiddleWare(adminUsecase, enums.LeaderAndMemberMiddleware)
	allAdminMiddleware := middlewares.NewAdminMiddleWare(adminUsecase, enums.AllAdminMiddleware)

	spectrumApi := r.Group("/spectrum")
	spectrumApi.POST("/create", leaderAndMemberMiddleware.AdminMiddleWare, handlers.CreateSpectrum)
	spectrumApi.PATCH("/update", leaderAndMemberMiddleware.AdminMiddleWare, handlers.UpdateSpectrum)
	spectrumApi.DELETE("/delete", leaderAndMemberMiddleware.AdminMiddleWare, handlers.DeleteSpectrum)
	spectrumApi.GET("", allAdminMiddleware.AdminMiddleWare, handlers.RetrieveSpectrums)
}
