package main

import (
	"espectro/internal/config"
	"espectro/internal/delivery/http/routes"
	"espectro/internal/infrastructure/database"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/gin-gonic/gin"
)

func main() {

	cfg, cfgErr := config.Load(".env")

	if cfgErr != nil {
		logger.Info("ENV file loading error : ", cfgErr)
		return
	}

	db, err := database.NewPostgres(cfg.DBUrl)

	if err != nil {
		logger.Error("Database connection error : ", err)
	}

	gin := gin.Default()

	routes.RegisterUserRoutes(gin)

	listenErr := gin.Run(":8080")

	if listenErr != nil {
		db.Close()
		logger.Error("Error while running the listener : ", listenErr)
	}

	logger.Info("Server started")
}
