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
	adminUsecases := usecases.NewAdminUsecases(adminRepo)
	galleryRepo := repositoryimple.NewGalleryPostgresRepo(db)
	venueRepo := repositoryimple.NewVenuePostgresRepo(db)
	mediaRepo := repositoryimple.NewMediaCloudinaryRepo(cld)
	eventRepo := repositoryimple.NewEventPostgresRepo(db)
	eventGalleryRepo := repositoryimple.NewEventGalleryPostgresRepo(db)
	transaction := database.NewTransactionManager(db)

	galleryUsecases := usecases.NewGalleryUsecases(galleryRepo, mediaRepo, venueRepo, eventGalleryRepo, eventRepo, transaction)

	handlers := handlers.NewGalleryHandlers(galleryUsecases)

	leaderAndMemberAdmin := middlewares.NewAdminMiddleWare(adminUsecases, enums.LeaderAndMemberMiddleware)

	galleryApi := r.Group("/gallery")

	galleryApi.POST("", leaderAndMemberAdmin.AdminMiddleWare, handlers.CreateGallery)
	galleryApi.PATCH("", leaderAndMemberAdmin.AdminMiddleWare, handlers.UpdateGallery)
	galleryApi.DELETE("", leaderAndMemberAdmin.AdminMiddleWare, handlers.DeleteGallery)
	galleryApi.GET("", handlers.RetrieveGalleries)
	galleryApi.POST("/add-event", leaderAndMemberAdmin.AdminMiddleWare, handlers.AddGalleryToEvent)
	galleryApi.DELETE("/from-event", leaderAndMemberAdmin.AdminMiddleWare, handlers.DeleteGalleryFromEvent)
}
