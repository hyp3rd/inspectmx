package logger

import (
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
	"io.hyperd/inspectmx/helpers"
)

var (
	once           sync.Once
	loggerInstance *ApplicationLogger
)

type ApplicationLogger struct {
	Logger *logrus.Logger
}

type WithFields map[string]interface{}

// instance creates a logger and makes it consumable throughout the application
func instance() *ApplicationLogger {
	once.Do(func() {
		loggerInstance = &ApplicationLogger{
			Logger: logrus.New(),
		}
		loggerInstance.Logger.AddHook(&DefaultFieldHook{GetValue: func() string { return "inspectmx" }})
	})
	return loggerInstance
}

func Fatal(msg string, fields map[string]interface{}) {
	instance().Logger.WithFields(fields).Fatal(msg)
}

func Info(msg string, fields map[string]interface{}) {
	instance().Logger.WithFields(fields).Info(msg)
}

func Error(msg string, fields map[string]interface{}) {
	instance().Logger.WithFields(fields).Error(msg)
}

func Warn(msg string, fields map[string]interface{}) {
	instance().Logger.WithFields(fields).Warn(msg)
}

func Debug(msg string, fields map[string]interface{}) {
	instance().Logger.WithFields(fields).Debug(msg)
}

func SetupLogger(logLevel string) {

	switch strings.ToLower(*helpers.String(logLevel)) {
	case "info", "warn":
		instance().Logger.SetFormatter(&logrus.TextFormatter{
			ForceColors:               true,
			DisableColors:             false,
			EnvironmentOverrideColors: true,
		})
	default:
		instance().Logger.SetFormatter(&logrus.JSONFormatter{
			PrettyPrint: true,
		})
		instance().Logger.SetLevel(logrus.DebugLevel)
		instance().Logger.Debug("log level is set to debugging")
	}
}
