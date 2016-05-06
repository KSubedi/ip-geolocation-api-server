package utils

import (
	"github.com/Sirupsen/logrus"
	"github.com/rifflock/lfshook"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Fill out the endpoint and uncomment if you want to log to them
	// endPoint := ""
	// sumoLogicHook := sumorus.NewSumoLogicHook(endPoint, "ip.api.fridev.com", logrus.InfoLevel, "ip_api")
	// _ = sumoLogicHook
	// logrus.AddHook(sumoLogicHook)

	lfsHook := lfshook.NewHook(lfshook.PathMap{
		logrus.InfoLevel:  "info.log",
		logrus.ErrorLevel: "error.log",
	})

	logrus.AddHook(lfsHook)
}

func LogInfo(i logrus.Fields, m string) {
	logrus.WithFields(i).Info(m)
}

func LogError(i logrus.Fields, m string) {
	logrus.WithFields(i).Error(m)
}
