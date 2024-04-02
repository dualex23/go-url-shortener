package config

import (
	"flag"
	"fmt"
	"os"
)

type App struct {
	BaseURL		string
	ServerAddr 	string
}

func AppParseFlags() *App {
	var appConfig App

	// значения по умолчанию.
    appConfig.ServerAddr = "localhost:8080"
    appConfig.BaseURL = "http://localhost:8080"
	
	flag.StringVar(&appConfig.ServerAddr, "a", appConfig.ServerAddr, "Адрес запуска HTTP-сервера")
    flag.StringVar(&appConfig.BaseURL, "b", appConfig.BaseURL, "Базовый адрес результирующего сокращённого URL")
    flag.Parse()

    if envAddr := os.Getenv("SERVER_ADDRESS"); envAddr != "" {
        appConfig.ServerAddr = envAddr
    }
    if envBaseURL := os.Getenv("BASE_URL"); envBaseURL != "" {
        appConfig.BaseURL = envBaseURL
    }

	fmt.Printf("APF Сервер: POST %s\n", appConfig.ServerAddr)
	fmt.Printf("APF Базовый адрес: GET %s\n", appConfig.BaseURL)

	return &appConfig
}