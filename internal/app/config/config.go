package config

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type App struct {
	BaseURL         string
	ServerAddr      string
	FileStoragePath string
}

func AppParseFlags() *App {
	var appConfig App

	baseDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	dataDir := filepath.Join(baseDir, "data", "tmp")

	appConfig.ServerAddr = "localhost:8080"
	appConfig.BaseURL = "http://localhost:8080"
	defaultFilePath := filepath.Join(dataDir, "short-url-db.json")

	flag.StringVar(&appConfig.ServerAddr, "a", appConfig.ServerAddr, "Адрес запуска HTTP-сервера")
	flag.StringVar(&appConfig.BaseURL, "b", appConfig.BaseURL, "Базовый адрес результирующего сокращённого URL")
	flag.StringVar(&appConfig.FileStoragePath, "f", defaultFilePath, "Путь к файлу данных")
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
	if envFilePath := os.Getenv("FILE_STORAGE_PATH"); envFilePath != "" {
		appConfig.FileStoragePath = envFilePath
	} else {
		appConfig.FileStoragePath = defaultFilePath
	}

	return &appConfig
}
