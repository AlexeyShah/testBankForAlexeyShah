package logger

import (
	"bankService/internal/helpers/consts"
	"os"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func init() {
	Logger = logrus.New()
	Logger.SetOutput(os.Stdout)

	if len(os.Getenv(consts.EnvironmentName)) == 0 || os.Getenv(consts.EnvironmentName) == consts.EnvironmentLocal || os.Getenv(consts.EnvironmentName) == consts.EnvironmentDebug {
		Logger.Info("log level = debug")
		Logger.SetLevel(logrus.DebugLevel)
	} else if os.Getenv(consts.EnvironmentName) == consts.EnvironmentTest || os.Getenv(consts.EnvironmentName) == consts.EnvironmentProduction {
		Logger.Info("log level = info")
		Logger.SetLevel(logrus.InfoLevel)
	} else {
		Logger.Info("log level = default")
	}
}

func GetLogger() *logrus.Logger {
	return Logger
}
