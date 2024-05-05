package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

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
	logger.GetLogger().Info("Starting server", zap.String("path", appConfig.FileStoragePath))

	if appConfig.FileStoragePath == "" {
		logger.GetLogger().Fatal("File storage path is not specified")
	}

	storageInstance := storage.NewStorage(appConfig.FileStoragePath)
	if storageInstance == nil {
		logger.GetLogger().Fatal("Failed to create storage object")
	}

	db, err := storage.NewDB(appConfig.DataBaseDSN)
	if err != nil {
		logger.GetLogger().Fatal("Failed to connect database", zap.Error(err))
	}
	defer db.Close()

	sh := handler.NewShortenerHandler(appConfig.BaseURL, storageInstance)

	r := chi.NewRouter()
	r.Use(middleware.GzipMiddleware, middleware.WithLogging)
	r.Post("/", sh.MainHandler)
	r.Get("/{id}", sh.GetHandler)
	r.Post("/api/shorten", sh.APIHandler)

	logger.GetLogger().Info("Server is started", zap.String("address", appConfig.ServerAddr))

	if err := http.ListenAndServe(appConfig.ServerAddr, r); err != nil {
		logger.GetLogger().Fatal("Server failed to start", zap.Error(err))
	}
}
