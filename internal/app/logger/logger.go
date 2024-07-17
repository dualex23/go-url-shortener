package logger

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Sugar *zap.SugaredLogger

func New() {
	config := zap.NewProductionConfig()

	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")

	logger, err := config.Build()
	if err != nil {
		log.Fatalf("Cannot initialize logger: %v", err)
	}

	Sugar = logger.Sugar()
}

func GetLogger() *zap.SugaredLogger {
	return Sugar
}
