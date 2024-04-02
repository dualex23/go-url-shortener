package handler

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

// var (
// 	baseURL = "http://localhost:8080"
// 	mapUrls = make(map[string]string) // "abcDeFgh: https://practicum.yandex.ru/"
// )

type ShortenerHandler struct {
	BaseURL string
	MapURLs map[string]string
}

func NewShortenerHandler(baseURL string) *ShortenerHandler {
	return &ShortenerHandler{
		BaseURL: baseURL,
		MapURLs: make(map[string]string),
	}
}

func (h *ShortenerHandler) MainHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST request is allowed!", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil || len(body) == 0 {
		http.Error(w, "Request body cannot be empty", http.StatusBadRequest)
		return
	}

	originalURL := string(body)

	id := uuid.New().String()[:8]

	h.MapURLs[id] = originalURL

	shortenedURL := id

	fmt.Printf("%s/%s", h.BaseURL, id)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("%s/%s", h.BaseURL, shortenedURL)))
}

func (h *ShortenerHandler) GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET request is allowed!", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/")
	if id == "" {
		http.Error(w, "Missing ID", http.StatusBadRequest)
		return
	}
	idiInMapUrls, ok := h.MapURLs[id]
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
	fmt.Fprintln(w, "Перенаправление на URL")
}