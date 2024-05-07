package handler

import (
	"encoding/json"
	"fmt"
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

	var requestData struct {
		URL string `json:"url"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	defer r.Body.Close()

	if err != nil {
		http.Error(w, "Error reading JSON", http.StatusBadRequest)
		return
	}

	if requestData.URL == "" {
		http.Error(w, "URL field is required", http.StatusBadRequest)
		return
	}

	_, id, err := h.Storage.Save(requestData.URL, h.BaseURL)
	if err != nil {
		http.Error(w, "Failed to create short URL", http.StatusInternalServerError)
		return
	}

	logger.GetLogger().Infoln(
		"handler:", "MainHandler",
		"method:", r.Method,
		"originalURL:", requestData.URL,
		"shortenedURL:", fmt.Sprintf("%s/%s", h.BaseURL, id),
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	response := map[string]string{id: fmt.Sprintf("%s/%s", h.BaseURL, id)}
	json.NewEncoder(w).Encode(response)
}

func (h *ShortenerHandler) GetHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/")

	fmt.Printf("GetHandler id = %s\n", id)

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

	logger.GetLogger().Infoln(
		"handler:", "GetHandler",
		"method:", r.Method,
		"requestUrl:", r.URL,
	)

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

	logger.GetLogger().Infoln(
		"handler", "APIHandler",
		"method", r.Method,
		"url", input.URL,
	)

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

	logger.GetLogger().Infoln(
		"response:", fmt.Sprintf("%s:%s", id, shortenedURL),
	)
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
