package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

var (
	baseURL = "http://localhost:8080"
	mapUrls = make(map[string]string) // "abcDeFgh: https://practicum.yandex.ru/"
)

func mainHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST request is allowed!", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "Could not read request body", http.StatusBadRequest)
		return
	}
	originalURL := string(body)

	id := uuid.New().String()[:8]
	mapUrls[id] = originalURL

	shortenedURL := id

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("%s/%s", baseURL, shortenedURL)))
}

func getHandler(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET request is allowed!", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/")
	if id == "" {
		http.Error(w, "Missing ID", http.StatusBadRequest)
		return
	}
	idiInMapUrls, ok := mapUrls[id]
	if !ok {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	returnedURL := idiInMapUrls

	if strings.HasPrefix(r.URL.Path, "/") && len(r.URL.Path) > 1 {
		w.Header().Set("Location", returnedURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		w.WriteHeader(http.StatusBadRequest)                   
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			mainHandler(w, r)
		case http.MethodGet:
			getHandler(w, r)
		default:
			http.Error(w, "Not Found", http.StatusBadRequest)
		}
	})

	err := http.ListenAndServe(":8080",mux)
	if err != nil {
		fmt.Print("Server hasn't been started")
		panic(err)
	}
}
