package config

import (
	"flag"
	"fmt"
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
	defaultFileName := "short-url-db"

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	fileName := defaultFileName
	flag.StringVar(&appConfig.ServerAddr, "a", appConfig.ServerAddr, "Адрес запуска HTTP-сервера")
	flag.StringVar(&appConfig.BaseURL, "b", appConfig.BaseURL, "Базовый адрес результирующего сокращённого URL")
	flag.StringVar(&fileName, "f", defaultFileName, "Базовое имя файла данных без расширения .json")
	flag.Parse()

	if envAddr := os.Getenv("SERVER_ADDRESS"); envAddr != "" {
		appConfig.ServerAddr = envAddr
	}
	if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
		appConfig.BaseURL = envBaseURL
	}

	if envFileName := os.Getenv("FILE_STORAGE_PATH"); envFileName != "" {
		fmt.Printf("Используется переменная окружения FILE_STORAGE_PATH=%s\n", envFileName)
		fileName = envFileName
	}

	appConfig.FileStoragePath = fmt.Sprintf("./../../tmp/%s.json", fileName)

	return &appConfig
}
