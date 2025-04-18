package main

import (
	docs "bankService/docs"
	"bankService/internal/handler"
	"bankService/internal/helpers/consts"
	"bankService/internal/logger"
	"bankService/internal/storage"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	logger.Logger.SetFormatter(&logrus.TextFormatter{})
}

// @title Api Bank Service
// @version 1.0
// @description bank rest application

// @host localhost:8080
// @BasePath /api/v1
func main() {
	logger.Logger.Info("Start application")

	err := storage.NewMigrator().Magrate()
	if err != nil {
		panic(err)
	}

	if os.Getenv(consts.EnvironmentName) == consts.EnvironmentProduction {
		gin.SetMode(gin.ReleaseMode)
	} else if os.Getenv(consts.EnvironmentName) == consts.EnvironmentTest {
		gin.SetMode(gin.TestMode)
	} else if os.Getenv(consts.EnvironmentName) == consts.EnvironmentDebug {
		gin.SetMode(gin.DebugMode)
	}

	router := gin.New()
	handler.RegisterRoutes(router)

	docs.SwaggerInfo.BasePath = "/api/v1"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	go func() {
		if err := router.Run(":8080"); err != nil {
			panic(err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan

	logger.Logger.Info("End application")
}
