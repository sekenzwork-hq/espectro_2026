package main

import (
	"espectro/config"
	"espectro/database"
	"espectro/routes"
	"espectro/services"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/gin-gonic/gin"
)

func main() {

	cfg, cfgErr := config.Load(".env")
	if cfgErr != nil {
		logger.Info("ENV file loading error : ", cfgErr)
		return
	}

	if len(cfg.Version) == 0 || len(cfg.Version) == 1 {
		logger.Info("Invalid server version")
		return
	}

	db, dbErr := database.NewPostgres(cfg.DBUrl)
	if dbErr != nil {
		logger.Error("Database connection error : ", dbErr)
		return
	}

	gormDB, ormErr := database.NewOrm(db)
	if ormErr != nil {
		logger.Info("ORM connection error : ", ormErr)
		return
	}

	defer db.Close()

	cld, cldErr := services.NewCloudinary(cfg.CloudinaryUrl)
	if cldErr != nil {
		logger.Info("Cloudinary establishing error : ", cldErr)
		return
	}

	redisClient := database.NewRedis(cfg.RedisIP)
	defer redisClient.Close()

	gin := gin.Default()

	apiVersion := "/api/" + cfg.Version
	api := gin.Group(apiVersion)

	routes.RegisterUserRoutes(api, gormDB)
	routes.RegisterAdminRoutes(api, gormDB)
	routes.RegisterSpectrumRoutes(api, gormDB, cld)
	routes.RegisterVenueRoutes(api, gormDB)
	routes.RegisterEventRoutes(api, gormDB)
	routes.RegisterPartnerRoutes(api, gormDB, cld)
	routes.RegisterInvestorRoutes(api, gormDB, cld)
	routes.RegisterSponsorRoutes(api, gormDB, cld)

	listenErr := gin.Run(":8080")
	if listenErr != nil {
		db.Close()
		logger.Error("Error while running the listener : ", listenErr)
	}

	logger.Info("Server started")
}
