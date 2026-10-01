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
	userRepo := repositoryimple.NewUserPostgresRepo(db)

	eventUsecases := usecases.NewEventUsecases(eventRepo, spectrumRepo, venueRepo, transactionManager, eventsGalleryRepo, eventsRegistrationRepo, userRegistrationRepo)
	adminUsecases := usecases.NewAdminUsecases(adminRepo)
	userUsecases := usecases.NewUserUsecases(userRepo)

	eventHandlers := handlers.NewEventHandlers(eventUsecases)

	leaderAndMemberMiddleware := middlewares.NewAdminMiddleWare(adminUsecases, enums.LeaderAndMemberMiddleware)
	userMiddleware := middlewares.NewUserMiddleware(userUsecases)
	userAdminMiddleware := middlewares.NewUserAdminMiddleware(userUsecases, adminUsecases)

	eventApi := r.Group("event")

	eventApi.POST("", leaderAndMemberMiddleware.AdminMiddleWare, eventHandlers.CreateEvent)
	eventApi.PATCH("", leaderAndMemberMiddleware.AdminMiddleWare, eventHandlers.UpdateEvent)
	eventApi.DELETE("", leaderAndMemberMiddleware.AdminMiddleWare, eventHandlers.DeleteEvent)
	eventApi.GET("", userAdminMiddleware.UserAdminMiddleware, eventHandlers.RetrieveEvents)

	registrationApi := eventApi.Group("/registration")

	registrationApi.POST("", userMiddleware.UserMiddleware, eventHandlers.Register)
	registrationApi.PATCH("/status", leaderAndMemberMiddleware.AdminMiddleWare, eventHandlers.ChangeRegistrationStatus)
	registrationApi.POST("/withdraw", userMiddleware.UserMiddleware, eventHandlers.WithdrawRegistration)
	registrationApi.GET("", userMiddleware.UserMiddleware, eventHandlers.RetrieveRegisteredEventsFromUserSide)
}
