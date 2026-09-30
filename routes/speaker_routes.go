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

func RegisterSpeakerRoutes(r *gin.RouterGroup, db *gorm.DB, cld *cloudinary.Cloudinary) {

	speakerRepo := repositoryimple.NewSpeakerPostgresRepo(db)
	adminRepo := repositoryimple.NewAdminPostgresRepo(db)
	eventRepo := repositoryimple.NewEventPostgresRepo(db)
	eventSpeakerRepo := repositoryimple.NewEventSpeakerPostgresRepo(db)
	mediaRepo := repositoryimple.NewMediaCloudinaryRepo(cld)

	speakerUsecases := usecases.NewSpeakerUsecases(speakerRepo, mediaRepo, eventRepo, eventSpeakerRepo)
	adminUsecases := usecases.NewAdminUsecases(adminRepo)

	speakerHandlers := handlers.NewSpeakerHandlers(speakerUsecases)

	leaderMemberAdminMiddleware := middlewares.NewAdminMiddleWare(adminUsecases, enums.LeaderAndMemberMiddleware)

	speakerApi := r.Group("/speaker")

	speakerApi.POST("", leaderMemberAdminMiddleware.AdminMiddleWare, speakerHandlers.CreateSpeaker)
	speakerApi.PATCH("", leaderMemberAdminMiddleware.AdminMiddleWare, speakerHandlers.UpdateSpeaker)
	speakerApi.DELETE("", leaderMemberAdminMiddleware.AdminMiddleWare, speakerHandlers.DeleteSpeaker)
	speakerApi.GET("", speakerHandlers.RetrieveSpeaker)

	eventSpeakerApi := speakerApi.Group("/event")

	eventSpeakerApi.POST("", leaderMemberAdminMiddleware.AdminMiddleWare, speakerHandlers.AddSpeakerToEvent)
	eventSpeakerApi.DELETE("", leaderMemberAdminMiddleware.AdminMiddleWare, speakerHandlers.RemoveSpeakerFromEvent)
}
