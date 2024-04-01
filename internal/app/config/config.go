package config

import (
	"flag"
	"fmt"
)

type App struct {
	BaseURL		string
	ServerAddr 	string
}

func AppParseFlags() *App {
	var serverAddr, baseURL string
	
	flag.StringVar(&serverAddr, "a", "localhost:8080", "Адрес запуска HTTP-сервера")
	flag.StringVar(&baseURL, "b", "http://localhost:8080", "Базовый адрес результирующего сокращённого URL")
	flag.Parse()

	fmt.Printf("Сервер: POST %s\n", serverAddr)
	fmt.Printf("Базовый адрес: GET %s\n", baseURL)

	return &App{
		ServerAddr: serverAddr,
		BaseURL:    baseURL,
	}
}