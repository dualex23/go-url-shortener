package storage

import (
	"fmt"

	"github.com/dualex23/go-url-shortener/internal/app/logger"
	"github.com/google/uuid"
)

func NewStorage(fileName string, mode string, db DataBaseInterface) *Storage {
	logger.GetLogger().Info("NewStorage")

	storage := &Storage{
		StorageMode: mode,
		StoragePath: fileName,
		DataBase:    db,
		UrlsMap:     make(map[string]URLData),
	}

	if err := storage.Load(); err != nil {
		logger.GetLogger().Errorf("Error loading storage: %v", err)
	}

	return storage
}

func (s *Storage) Load() error {
	logger.GetLogger().Info("Load")

	switch s.StorageMode {
	case "db":
		urls, err := s.DataBase.LoadUrls()
		if err != nil {
			logger.GetLogger().Errorf("Failed to load URLs from database: %v", err)
			return err
		}
		s.UrlsMap = urls
	case "file":
		if s.StoragePath == "" {
			logger.GetLogger().Error("No file path specified for file-based storage.")
			return nil
		}

		return s.LoadUrlFromFile()
	default:
		logger.GetLogger().Info("Using in-memory storage, no initial data loading required.")
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
