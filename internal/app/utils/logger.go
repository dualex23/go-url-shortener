package utils

import (
	"log"

	"go.uber.org/zap"
)

var Sugar *zap.SugaredLogger

func InitLogger() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Cannot initialize logger: %v", err)
	}

	Sugar = logger.Sugar()
}

func GetLogger() *zap.SugaredLogger {
	return Sugar
}
