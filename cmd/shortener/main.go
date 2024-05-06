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

	var storageMode string
	if appConfig.DataBaseDSN != "" {
		storageMode = "db"
	} else if appConfig.FileStoragePath != "" {
		storageMode = "file"
	} else {
		storageMode = "memory"
	}

	logger.GetLogger().Infof("mode=%s", storageMode)

	if appConfig.FileStoragePath == "" {
		logger.GetLogger().Fatal("File storage path is not specified\n")
	}

	db, err := storage.NewDB(appConfig.DataBaseDSN)
	if err != nil {
		logger.GetLogger().Fatal("Failed to connect database\n", zap.Error(err))
	}
	defer db.Close()

	storageInstance := storage.NewStorage(appConfig.FileStoragePath, storageMode, db)
	if storageInstance == nil {
		logger.GetLogger().Fatal("Failed to create storage object\n")
	}

	sh := handler.NewShortenerHandler(appConfig.BaseURL, storageInstance)

	r := chi.NewRouter()
	r.Use(middleware.GzipMiddleware, middleware.WithLogging)
	r.Post("/", sh.MainHandler)
	r.Get("/{id}", sh.GetHandler)
	r.Post("/api/shorten", sh.APIHandler)
	r.Get("/ping", sh.PingTest)

	logger.GetLogger().Info("Server is started\n", zap.String("address", appConfig.ServerAddr))

	if err := http.ListenAndServe(appConfig.ServerAddr, r); err != nil {
		logger.GetLogger().Fatal("Server failed to start", zap.Error(err))
	}
}
