package routes

import (
	"espectro/database"
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
	transactionM := database.NewTransactionManager(db)
	cldMediaRepo := repositoryimple.NewMediaCloudinaryRepo(cld)
	eventRepo := repositoryimple.NewEventPostgresRepo(db)
	userRepo := repositoryimple.NewUserPostgresRepo(db)

	spectrumUsecase := usecases.NewSpectrumUsecases(spectrumRepo, cldMediaRepo, eventRepo, transactionM)
	adminUsecases := usecases.NewAdminUsecases(adminRepo)
	userUsecases := usecases.NewUserUsecases(userRepo)

	handlers := handlers.NewSpectrumHandlers(spectrumUsecase)

	leaderAndMemberMiddleware := middlewares.NewAdminMiddleWare(adminUsecases, enums.LeaderAndMemberMiddleware)
	userAdminMiddleware := middlewares.NewUserAdminMiddleware(userUsecases, adminUsecases)

	spectrumApi := r.Group("/spectrum")

	spectrumApi.POST("", leaderAndMemberMiddleware.AdminMiddleWare, handlers.CreateSpectrum)
	spectrumApi.PATCH("", leaderAndMemberMiddleware.AdminMiddleWare, handlers.UpdateSpectrum)
	spectrumApi.DELETE("", leaderAndMemberMiddleware.AdminMiddleWare, handlers.DeleteSpectrum)
	spectrumApi.GET("", userAdminMiddleware.UserAdminMiddleware, handlers.RetrieveSpectrums)
}
