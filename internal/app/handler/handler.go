package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/dualex23/go-url-shortener/internal/app/storage"
	"github.com/google/uuid"
)

type ShortenerHandler struct {
	BaseURL string
	Storage *storage.Storage
	mx      sync.RWMutex
}

func NewShortenerHandler(baseURL string, storage *storage.Storage) *ShortenerHandler {
	return &ShortenerHandler{
		BaseURL: baseURL,
		Storage: storage,
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

	urlData := storage.URLData{
		ID:          id,
		OriginalURL: originalURL,
		ShortURL:    fmt.Sprintf("%s/%s", h.BaseURL, id),
	}

	h.Storage.UrlsData = append(h.Storage.UrlsData, urlData)
	if err := h.Storage.SaveURLsData(); err != nil {
		http.Error(w, "Failed to save data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(urlData.ShortURL))
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

	var originalURL string
	found := false
	for _, data := range h.Storage.UrlsData {
		if data.ID == id {
			originalURL = data.OriginalURL
			found = true
			break
		}
	}

	if !found {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	if strings.HasPrefix(r.URL.Path, "/") && len(r.URL.Path) > 1 {
		w.Header().Set("Location", originalURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (h *ShortenerHandler) APIHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST request is allowed!", http.StatusMethodNotAllowed)
		return
	}

	var input struct {
		URL string `json:"url"`
	}
	err := json.NewDecoder(r.Body).Decode(&input)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, "Error reading JSON", http.StatusBadRequest)
		return
	}

	if input.URL == "" {
		http.Error(w, "URL field is required", http.StatusBadRequest)
		return
	}

	id := uuid.New().String()[:8]
	shortenedURL := fmt.Sprintf("%s/%s", h.BaseURL, id)

	urlData := storage.URLData{
		ID:          id,
		OriginalURL: input.URL,
		ShortURL:    shortenedURL,
	}

	h.Storage.UrlsData = append(h.Storage.UrlsData, urlData)
	if err := h.Storage.SaveURLsData(); err != nil {
		http.Error(w, "Failed to save data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"result": shortenedURL})
}
