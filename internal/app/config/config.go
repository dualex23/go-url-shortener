package config

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/dualex23/go-url-shortener/internal/app/logger"
	"github.com/joho/godotenv"
)

type App struct {
	BaseURL         string
	ServerAddr      string
	FileStoragePath string
	DataBaseDSN     string
}

func AppParseFlags() *App {
	var appConfig App

	appConfig.ServerAddr = "localhost:8080"
	appConfig.BaseURL = "http://localhost:8080"
	defaultFilePath := "/tmp/short-url-db.json"

	flag.StringVar(&appConfig.ServerAddr, "a", appConfig.ServerAddr, "Адрес запуска HTTP-сервера")
	flag.StringVar(&appConfig.BaseURL, "b", appConfig.BaseURL, "Базовый адрес результирующего сокращённого URL")
	flag.StringVar(&appConfig.FileStoragePath, "f", defaultFilePath, "Имя файла данных без пути")
	flag.StringVar(&appConfig.DataBaseDSN, "d", appConfig.DataBaseDSN, "DB настройки")
	flag.Parse()

	err := godotenv.Load(".env")
	if err != nil {
		logger.GetLogger().Error("Warning: .env file not found or error loading .env file")
	}

	if envAddr := os.Getenv("SERVER_ADDRESS"); envAddr != "" {
		logger.GetLogger().Infof("env SERVER_ADDRESS = %s", envAddr)
		appConfig.ServerAddr = envAddr
	}
	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		logger.GetLogger().Infof("env BASE_URL = %s", envBaseURL)
		appConfig.BaseURL = envBaseURL
	}
	if envDatabaseDSN := os.Getenv("DATABASE_DSN"); envDatabaseDSN != "" && appConfig.DataBaseDSN == "" {
		appConfig.DataBaseDSN = envDatabaseDSN
	}

	currentDir, _ := os.Getwd()
	appConfig.FileStoragePath = filepath.Join(currentDir, appConfig.FileStoragePath)

	if envFilePath := os.Getenv("FILE_STORAGE_PATH"); envFilePath != "" {
		appConfig.FileStoragePath = envFilePath
	}

	return &appConfig
}
