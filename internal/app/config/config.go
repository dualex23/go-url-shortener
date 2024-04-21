package config

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
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
	defaultFileName := "./tmp/short-url-db.json"

	flag.StringVar(&appConfig.ServerAddr, "a", appConfig.ServerAddr, "Адрес запуска HTTP-сервера")
	flag.StringVar(&appConfig.BaseURL, "b", appConfig.BaseURL, "Базовый адрес результирующего сокращённого URL")
	flag.StringVar(&appConfig.FileStoragePath, "f", defaultFileName, "Базовое имя файла данных без расширения .json")
	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found or error loading .env file")
	}

	if envAddr := os.Getenv("SERVER_ADDRESS"); envAddr != "" {
		appConfig.ServerAddr = envAddr
	}
	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		appConfig.BaseURL = envBaseURL
	}

	if envFileName := os.Getenv("FILE_STORAGE_PATH"); envFileName != "" {
		appConfig.FileStoragePath = envFileName
	}

	return &appConfig
}
