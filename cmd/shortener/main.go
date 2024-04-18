package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/dualex23/go-url-shortener/internal/app/config"
	"github.com/dualex23/go-url-shortener/internal/app/handler"
	"github.com/dualex23/go-url-shortener/internal/app/middleware"
	"github.com/dualex23/go-url-shortener/internal/app/utils"
)

func main() {
	utils.InitLogger()
	defer utils.GetLogger().Sync()

	appConfig := config.AppParseFlags()

	r := chi.NewRouter()

	shortenerHandler := handler.NewShortenerHandler(appConfig.BaseURL)

	r.Use(middleware.WithLogging)
	r.Post("/", shortenerHandler.MainHandler)
	r.Get("/{id}", shortenerHandler.GetHandler)

	err := http.ListenAndServe(appConfig.ServerAddr, r)
	if err != nil {
		log.Fatal(err)
	}
}
