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
	"github.com/lpernett/godotenv"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	logger.Logger.SetFormatter(&logrus.TextFormatter{})
	err := godotenv.Load("config.env")
	if err != nil {
		logger.Logger.Fatalf("failed load env: %v", err)
	}

	logger.Logger.Infof("Environment = %s", os.Getenv(consts.EnvironmentName))
	logger.Logger.Infof("connection = %s", os.Getenv(consts.CloudConnectionPostgre))
	if len(os.Getenv(consts.EnvironmentName)) == 0 || os.Getenv(consts.EnvironmentName) == consts.EnvironmentLocal || os.Getenv(consts.EnvironmentName) == consts.EnvironmentDebug {
		logger.Logger.Info("log level = debug")
		logger.Logger.SetLevel(logrus.DebugLevel)
	} else if os.Getenv(consts.EnvironmentName) == consts.EnvironmentTest || os.Getenv(consts.EnvironmentName) == consts.EnvironmentProduction {
		logger.Logger.Info("log level = info")
		logger.Logger.SetLevel(logrus.InfoLevel)
	} else {
		logger.Logger.Info("log level = default")
	}
}

// @title Api Bank Service
// @version 1.0
// @description bank rest application

// @host localhost:8080
// @BasePath /api/v1
func main() {
	logger.Logger.Info("Start application")

	conn, err := storage.NewPostgreConnection()
	if err != nil {
		panic(err)
	}
	err = storage.NewMigrator(conn).Magrate()
	if err != nil {
		panic(err)
	}

	logger.Logger.Info("migration success")

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
