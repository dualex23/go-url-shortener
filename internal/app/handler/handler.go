package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/dualex23/go-url-shortener/internal/app/logger"
	"github.com/dualex23/go-url-shortener/internal/app/storage"
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

	storageSave, id, err := h.Storage.Save(originalURL, h.BaseURL)
	if err != nil {
		http.Error(w, "Failed to create short URL", http.StatusInternalServerError)
		return
	}
	fmt.Printf("test %v", storageSave)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("http://localhost:8080/%s", id)))

}

func (h *ShortenerHandler) GetHandler(w http.ResponseWriter, r *http.Request) {

	id := strings.TrimPrefix(r.URL.Path, "/")

	logger.GetLogger().Infof("id=%s\n", id)

	if id == "" {
		http.Error(w, "Missing ID", http.StatusBadRequest)
		return
	}

	originalURL, err := h.Storage.FindByID(id)
	if err != nil {
		logger.GetLogger().Errorf("Error finding URL: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)

		return
	}

	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
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

	shortenedURL, id, err := h.Storage.Save(input.URL, h.BaseURL)
	if err != nil {
		http.Error(w, "Failed to create or save URL", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{id: shortenedURL})
}

func (h *ShortenerHandler) PingTest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only Get request is allowed!", http.StatusMethodNotAllowed)
		return
	}

	err := h.Storage.DataBase.Ping()
	if err != nil {
		logger.GetLogger().Errorf("Database connection failed: %v", err)

		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		h.Storage.DataBase.Close()
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Database connection successful"))
}
