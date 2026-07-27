package routes

import (
	"espectro/database"
	"espectro/enums"
	"espectro/handlers"
	"espectro/middlewares"
	repositoryimple "espectro/repository_imple"
	"espectro/usecases"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterEventRoutes(r *gin.RouterGroup, db *gorm.DB) {

	eventRepo := repositoryimple.NewEventPostgresRepo(db)
	venueRepo := repositoryimple.NewVenuePostgresRepo(db)
	spectrumRepo := repositoryimple.NewSpectrumPostgresRepo(db)
	adminRepo := repositoryimple.NewAdminPostgresRepo(db)
	transactionManager := database.NewTransactionManager(db)
	eventsGalleryRepo := repositoryimple.NewEventGalleryPostgresRepo(db)
	eventsRegistrationRepo := repositoryimple.NewEventRegistrationPostgresRepo(db)
	userRegistrationRepo := repositoryimple.NewUserPostgresRepo(db)
	eventUsecases := usecases.NewEventUsecases(eventRepo, spectrumRepo, venueRepo, transactionManager, eventsGalleryRepo, eventsRegistrationRepo, userRegistrationRepo)
	adminUsecases := usecases.NewAdminUsecases(adminRepo)

	eventHandlers := handlers.NewEventHandlers(eventUsecases)

	leaderAndMemberMiddleware := middlewares.NewAdminMiddleWare(adminUsecases, enums.LeaderAndMemberMiddleware)

	eventApi := r.Group("event")

	eventApi.POST("/create", leaderAndMemberMiddleware.AdminMiddleWare, eventHandlers.CreateEvent)
	eventApi.PATCH("/update", leaderAndMemberMiddleware.AdminMiddleWare, eventHandlers.UpdateEvent)
	eventApi.DELETE("/delete", leaderAndMemberMiddleware.AdminMiddleWare, eventHandlers.DeleteEvent)
	eventApi.GET("", eventHandlers.RetrieveEvents)

	registrationApi := eventApi.Group("/registration")

	registrationApi.POST("/register", eventHandlers.Register)
	registrationApi.PATCH("/change-status", leaderAndMemberMiddleware.AdminMiddleWare, eventHandlers.ChangeRegistrationStatus)
	registrationApi.POST("/withdraw", eventHandlers.WithdrawRegistration)
}
