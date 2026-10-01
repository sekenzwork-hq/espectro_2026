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

func RegisterGalleryRoutes(r *gin.RouterGroup, db *gorm.DB, cld *cloudinary.Cloudinary) {

	adminRepo := repositoryimple.NewAdminPostgresRepo(db)
	galleryRepo := repositoryimple.NewGalleryPostgresRepo(db)
	venueRepo := repositoryimple.NewVenuePostgresRepo(db)
	mediaRepo := repositoryimple.NewMediaCloudinaryRepo(cld)
	eventRepo := repositoryimple.NewEventPostgresRepo(db)
	eventGalleryRepo := repositoryimple.NewEventGalleryPostgresRepo(db)
	transaction := database.NewTransactionManager(db)
	userRepo := repositoryimple.NewUserPostgresRepo(db)

	galleryUsecases := usecases.NewGalleryUsecases(galleryRepo, mediaRepo, venueRepo, eventGalleryRepo, eventRepo, transaction)
	adminUsecases := usecases.NewAdminUsecases(adminRepo)
	userUsecases := usecases.NewUserUsecases(userRepo)

	handlers := handlers.NewGalleryHandlers(galleryUsecases)

	leaderAndMemberAdmin := middlewares.NewAdminMiddleWare(adminUsecases, enums.LeaderAndMemberMiddleware)
	userAdminMiddleware := middlewares.NewUserAdminMiddleware(userUsecases, adminUsecases)

	galleryApi := r.Group("/gallery")

	galleryApi.POST("", leaderAndMemberAdmin.AdminMiddleWare, handlers.CreateGallery)
	galleryApi.PATCH("", leaderAndMemberAdmin.AdminMiddleWare, handlers.UpdateGallery)
	galleryApi.DELETE("", leaderAndMemberAdmin.AdminMiddleWare, handlers.DeleteGallery)
	galleryApi.GET("", userAdminMiddleware.UserAdminMiddleware, handlers.RetrieveGalleries)
	galleryApi.POST("/add-event", leaderAndMemberAdmin.AdminMiddleWare, handlers.AddGalleryToEvent)
	galleryApi.DELETE("/from-event", leaderAndMemberAdmin.AdminMiddleWare, handlers.DeleteGalleryFromEvent)
}
