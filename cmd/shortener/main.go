package main

import (
	// стандартные библиотеки

	// сторонние библиотеки
	"net/http"

	"github.com/go-chi/chi/v5"

	// дергаем свои пакеты
	"github.com/dualex23/go-url-shortener/internal/app/config"
	"github.com/dualex23/go-url-shortener/internal/app/handlers"
)

func main() {
	c := &config.App{
		Port: ":8080",
	}

	r := chi.NewRouter()

	r.Post("/", handlers.MainHandler)
	r.Get("/{id}",handlers.GetHandler) 

	err := http.ListenAndServe(c.Port, r)
	if err != nil {
		panic(err)
	}
}
