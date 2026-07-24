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

func RegisterGalleryRoutes(r *gin.RouterGroup, db *gorm.DB, cld *cloudinary.Cloudinary) {

	adminRepo := repositoryimple.NewAdminPostgresRepo(db)
	adminUsecases := usecases.NewAdminUsecases(adminRepo)
	galleryRepo := repositoryimple.NewGalleryPostgresRepo(db)
	venueRepo := repositoryimple.NewVenuePostgresRepo(db)
	mediaRepo := repositoryimple.NewMediaCloudinaryRepo(cld)
	eventRepo := repositoryimple.NewEventPostgresRepo(db)
	eventGalleryRepo := repositoryimple.NewEventGalleryPostgresRepo(db)
	galleryUsecases := usecases.NewGalleryUsecases(galleryRepo, mediaRepo, venueRepo, eventGalleryRepo, eventRepo)
	handlers := handlers.NewGalleryHandlers(galleryUsecases)
	adminMiddleware := middlewares.NewAdminMiddleWare(adminUsecases, enums.LeaderAndMemberMiddleware)

	galleryApi := r.Group("/gallery")

	galleryApi.POST("/create", adminMiddleware.AdminMiddleWare, handlers.CreateGallery)
	galleryApi.PATCH("/update", adminMiddleware.AdminMiddleWare, handlers.UpdateGallery)
	galleryApi.DELETE("/delete", adminMiddleware.AdminMiddleWare, handlers.DeleteGallery)
	galleryApi.GET("", handlers.RetrieveGalleries)
	galleryApi.POST("/add-event", adminMiddleware.AdminMiddleWare, handlers.AddGalleryToEvent)
	galleryApi.DELETE("/delete-from-event", adminMiddleware.AdminMiddleWare, handlers.DeleteGalleryFromEvent)
}
