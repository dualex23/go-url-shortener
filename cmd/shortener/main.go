package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/dualex23/go-url-shortener/internal/app/config"
	"github.com/dualex23/go-url-shortener/internal/app/handler"
	"github.com/dualex23/go-url-shortener/internal/app/middleware"
	"github.com/dualex23/go-url-shortener/internal/app/storage"
	"github.com/dualex23/go-url-shortener/internal/app/utils"
)

func main() {
	utils.InitLogger()
	defer utils.GetLogger().Sync()

	appConfig := config.AppParseFlags()
	fmt.Printf("main FileStoragePath = %v\n", appConfig.FileStoragePath)
	storage.Init(appConfig.FileStoragePath)

	r := chi.NewRouter()

	shortenerHandler := handler.NewShortenerHandler(appConfig.BaseURL)

	r.Use(middleware.GzipMiddleware, middleware.WithLogging)
	r.Post("/", shortenerHandler.MainHandler)
	r.Get("/{id}", shortenerHandler.GetHandler)
	r.Post("/api/shorten", shortenerHandler.APIHandler)

	fmt.Printf("Server is started: %s\n", appConfig.ServerAddr)
	err := http.ListenAndServe(appConfig.ServerAddr, r)
	if err != nil {
		log.Fatal(err)
	}
}
