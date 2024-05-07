package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/dualex23/go-url-shortener/internal/app/logger"
	"github.com/dualex23/go-url-shortener/internal/app/storage"
)

type ShortenerHandler struct {
	BaseURL string
	Storage *storage.Storage
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

	// написать проверку на существующий url
	existingID, _, err := h.Storage.DataBase.FindByOriginalURL(originalURL)
	if err == nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(fmt.Sprintf("%s/%s", h.BaseURL, existingID)))
		return
	}

	_, id, err := h.Storage.Save(originalURL, h.BaseURL)
	if err != nil {
		http.Error(w, "Failed to create short URL", http.StatusInternalServerError)
		return
	}

	logger.GetLogger().Infoln(
		"handler", "APIHandler",
		"method", r.Method,
		"originalURL", originalURL,
		"BaseURL", h.BaseURL,
	)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("%s/%s", h.BaseURL, id)))

}

func (h *ShortenerHandler) GetHandler(w http.ResponseWriter, r *http.Request) {
	logger.GetLogger().Infoln(
		"method:", r.Method,
		"requestUrl:", r.URL,
		"fullPath:", r.URL.Host,
	)

	id := strings.TrimPrefix(r.URL.Path, "/")

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
		"method:", r.Method,
		"originalURL:", originalURL,
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

	_, existingShortened, err := h.Storage.DataBase.FindByOriginalURL(input.URL)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{"result": existingShortened})
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
	json.NewEncoder(w).Encode(map[string]string{"result": shortenedURL})
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

func (h *ShortenerHandler) BatchShortenHandler(w http.ResponseWriter, r *http.Request) {
	var req []storage.BatchShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var resp []storage.BatchShortenResponse
	var batch []storage.URLData

	for _, item := range req {
		id := storage.GenerateID()
		shortURL := fmt.Sprintf("%s/%s", h.BaseURL, id)

		batch = append(batch, storage.URLData{
			ID:          id,
			OriginalURL: item.OriginalURL,
			ShortURL:    shortURL,
		})

		resp = append(resp, storage.BatchShortenResponse{
			CorrelationID: item.CorrelationID,
			ShortURL:      shortURL,
		})
	}

	// Сохранение всех URL-ов в одной транзакции
	if err := h.Storage.DataBase.BatchSaveUrls(batch); err != nil {
		http.Error(w, "Failed to save batch URLs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
