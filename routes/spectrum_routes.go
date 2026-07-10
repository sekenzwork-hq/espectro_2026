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

func RegisterSpectrumRoutes(r *gin.RouterGroup, db *gorm.DB) {

	spectrumRepo := repositoryimple.NewSpectrumPostgresRepo(db)
	adminRepo := repositoryimple.NewAdminPostgresRepo(db)
	spectrumUsecase := usecases.NewSpectrumUsecases(spectrumRepo)
	adminUsecase := usecases.NewAdminUsecases(adminRepo)
	handlers := handlers.NewSpectrumHandlers(spectrumUsecase)

	leaderAndMemberMiddleware := middlewares.NewAdminMiddleWare(adminUsecase, enums.LeaderAndMemberMiddleware)

	spectrumApi := r.Group("/spectrum")
	spectrumApi.POST("/create", leaderAndMemberMiddleware.AdminMiddleWare, handlers.CreateSpectrum)

}
