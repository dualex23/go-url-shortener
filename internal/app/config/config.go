package config

import (
	"flag"
	"os"
)

type App struct {
	BaseURL         string
	ServerAddr      string
	FileStoragePath string
}

func AppParseFlags() *App {
	var appConfig App

	appConfig.ServerAddr = "localhost:8080"
	appConfig.BaseURL = "http://localhost:8080"
	appConfig.FileStoragePath = "./tmp/short-url-db.json"

	flag.StringVar(&appConfig.ServerAddr, "a", appConfig.ServerAddr, "Адрес запуска HTTP-сервера")
	flag.StringVar(&appConfig.BaseURL, "b", appConfig.BaseURL, "Базовый адрес результирующего сокращённого URL")
	flag.StringVar(&appConfig.FileStoragePath, "f", appConfig.FileStoragePath, "Адрес хранилища данных")
	flag.Parse()

	if envAddr := os.Getenv("SERVER_ADDRESS"); envAddr != "" {
		appConfig.ServerAddr = envAddr
	}
	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		appConfig.BaseURL = envBaseURL
	}
	if storagePath := os.Getenv("FILE_STORAGE_PATH"); storagePath != "" {
		appConfig.FileStoragePath = storagePath
	}

	return &appConfig
}
