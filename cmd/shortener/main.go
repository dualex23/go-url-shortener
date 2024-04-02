package main

import (
	// стандартные библиотеки

	// сторонние библиотеки

	"net/http"

	"github.com/go-chi/chi/v5"

	// дергаем свои пакеты

	"github.com/dualex23/go-url-shortener/internal/app/config"
	"github.com/dualex23/go-url-shortener/internal/app/handler"
)



func main() {
	appConfig := config.AppParseFlags()
	
	r := chi.NewRouter()
	
	shortenerHandler := handler.NewShortenerHandler(appConfig.BaseURL)

	r.Post("/", shortenerHandler.MainHandler)
	r.Get("/{id}",shortenerHandler.GetHandler)

	// debug
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Привет "+appConfig.ServerAddr +" "+appConfig.BaseURL))
	})
	
	err := http.ListenAndServe(appConfig.ServerAddr, r)
	if err != nil {
		panic(err)
	}
}
