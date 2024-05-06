package storage

import (
	"fmt"

	"github.com/dualex23/go-url-shortener/internal/app/logger"
	"github.com/google/uuid"
)

func NewStorage(fileName string, mode string, db DataBaseInterface) *Storage {
	logger.GetLogger().Info("NewStorage\n")

	storage := &Storage{
		StorageMode: mode,
		StoragePath: fileName,
		DataBase:    db,
		UrlsMap:     make(map[string]URLData),
	}

	if err := storage.Load(); err != nil {
		logger.GetLogger().Errorf("Error loading storage: %v\n", err)
	}

	return storage
}

func (s *Storage) Load() error {
	logger.GetLogger().Info("Load\n")

	switch s.StorageMode {
	case "db":
		urls, err := s.DataBase.LoadUrls()
		if err != nil {
			logger.GetLogger().Errorf("Failed to load URLs from database: %v\n", err)
			return err
		}
		s.UrlsMap = urls
	case "file", "memory":
		if s.StoragePath == "" {
			logger.GetLogger().Error("No file path specified for file-based storage.\n")
			return nil
		}

		return s.LoadUrlFromFile()
	default:
		logger.GetLogger().Info("Load: Using in-memory storage, no initial data loading required.\n")
	}

	return nil
}

func (s *Storage) Save(originalURL string, baseURL string) (string, string, error) {
	id := uuid.New().String()[:8]
	shortURL := fmt.Sprintf("%s/%s", baseURL, id)

	urlData := URLData{
		ID:          id,
		OriginalURL: originalURL,
		ShortURL:    shortURL,
	}

	switch s.StorageMode {
	case "db":
		if err := s.DataBase.SaveUrls(urlData.ID, urlData.ShortURL, urlData.OriginalURL); err != nil {
			return "", "", err
		}
	case "file":
		s.mu.Lock()
		s.UrlsMap[id] = urlData
		s.mu.Unlock()

		if err := s.SaveUrlToFile(); err != nil {
			return "", "", err
		}
	default:
		s.mu.Lock()
		s.UrlsMap[id] = urlData
		s.mu.Unlock()
	}

	return shortURL, id, nil
}

func (s *Storage) FindByID(id string) (string, error) {
	switch s.StorageMode {
	case "db":
		urlData, err := s.DataBase.LoadUrlByID(id)
		if err != nil {
			logger.GetLogger().Errorf("Failed to load URL from database by ID %s: %v", id, err)
			return "", err
		}
		return urlData.OriginalURL, nil
	case "file", "memory":
		urlData, ok := s.UrlsMap[id]
		if !ok {
			logger.GetLogger().Errorf("No URL found with ID %s", id)
			return "", fmt.Errorf("no URL found with ID %s", id)
		}
		return urlData.OriginalURL, nil
	default:
		logger.GetLogger().Errorf("Unknown storage mode")
		return "", fmt.Errorf("unknown storage mode")
	}
}
