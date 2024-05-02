package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/dualex23/go-url-shortener/internal/app/config"
	"github.com/dualex23/go-url-shortener/internal/app/handler"
	"github.com/dualex23/go-url-shortener/internal/app/logger"
	"github.com/dualex23/go-url-shortener/internal/app/middleware"
	"github.com/dualex23/go-url-shortener/internal/app/storage"
)

func main() {
	logger.New()
	defer logger.GetLogger().Sync()

	appConfig := config.AppParseFlags()
	fmt.Printf("main FileStoragePath = %v\n", appConfig.FileStoragePath)

	if appConfig.FileStoragePath == "" {
		log.Fatal("Не указан путь к файлу хранилища")
	}

	storage := storage.NewStorage(appConfig.FileStoragePath)
	if storage == nil {
		log.Fatal("Не удалось создать объект хранилища")
	}

	sh := handler.NewShortenerHandler(appConfig.BaseURL, storage)

	r := chi.NewRouter()

	r.Use(middleware.GzipMiddleware, middleware.WithLogging)
	r.Post("/", sh.MainHandler)
	r.Get("/{id}", sh.GetHandler)
	r.Post("/api/shorten", sh.APIHandler)

	fmt.Printf("Server is started: %s\n", appConfig.ServerAddr)
	err := http.ListenAndServe(appConfig.ServerAddr, r)
	if err != nil {
		log.Fatal(err)
	}
}
