package storage

import (
	"fmt"

	"github.com/dualex23/go-url-shortener/internal/app/logger"
	"github.com/google/uuid"
)

func GenerateID() string {
	id := uuid.New().String()[:8]
	return id
}

func NewStorage(fileName string, mode string, db DataBaseInterface) *Storage {
	//logger.GetLogger().Info("NewStorage\n")

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
	//logger.GetLogger().Info("Load\n")

	switch s.StorageMode {
	case "db":
		urls, err := s.DataBase.LoadUrls()
		if err != nil {
			logger.GetLogger().Errorf("Failed to load URLs from database: %v\n", err)
			return err
		}
		s.UrlsMap = urls
	case "file":
		if s.StoragePath == "" {
			logger.GetLogger().Error("No file path specified for file-based storage.\n")
			return nil
		}

		return s.LoadURLFromFile()
	default:
		fmt.Printf("s.UrlsMap=%v\n", s.UrlsMap)
		logger.GetLogger().Info("Load: Using in-memory storage, no initial data loading required.\n")
	}

	return nil
}

func (s *Storage) Save(originalURL string, baseURL string) (string, string, error) {
	logger.GetLogger().Info("Storage Save\n")

	id := GenerateID()
	shortURL := fmt.Sprintf("%s/%s", baseURL, id)

	logger.GetLogger().Infoln(
		"id:", id,
		"shortURL:", shortURL,
	)

	urlData := URLData{
		ID:          id,
		OriginalURL: originalURL,
		ShortURL:    shortURL,
	}

	switch s.StorageMode {
	case "db":
		fmt.Printf("StorageMode db\n")
		if err := s.DataBase.SaveUrls(urlData.ID, urlData.ShortURL, urlData.OriginalURL); err != nil {
			return "", "", err
		}
	case "file":
		fmt.Printf("StorageMode file\n")
		s.mu.Lock()
		s.UrlsMap[id] = urlData
		s.mu.Unlock()

		if err := s.SaveURLToFile(); err != nil {
			return "", "", err
		}
	case "memory":
		fmt.Printf("StorageMode memory\n")
		s.mu.Lock()
		s.UrlsMap[id] = urlData
		s.mu.Unlock()

		//fmt.Printf("s.UrlsMap=%v\n", s.UrlsMap)
		fmt.Printf("urlData=%v\n", urlData)
	default:
		fmt.Printf("StorageMode default\n")
		s.mu.Lock()
		s.UrlsMap[id] = urlData
		s.mu.Unlock()
	}

	logger.GetLogger().Infoln(
		"Name:", "Save return",
		"shortURL:", shortURL,
		"id:", id,
	)

	return shortURL, id, nil
}

func (s *Storage) FindByID(id string) (string, error) {
	logger.GetLogger().Infof("FindByID - %s", id)

	switch s.StorageMode {
	case "db":
		urlData, err := s.DataBase.LoadURLByID(id)
		if err != nil {
			logger.GetLogger().Errorf("Failed to load URL from database by ID %s: %v", id, err)
			return "", err
		}
		return urlData.OriginalURL, nil
	case "file", "memory":
		s.mu.Lock()
		urlData, ok := s.UrlsMap[id]
		s.mu.Unlock()

		if !ok {
			logger.GetLogger().Errorf("No URL found with ID %s", id)
			return "", fmt.Errorf("no URL found with ID %s", id)
		}
		logger.GetLogger().Infoln(
			"ID:", urlData.ID,
			"ShortUrl:", urlData.ShortURL,
			"OriginalUrl", urlData.OriginalURL,
		)
		return urlData.OriginalURL, nil
	default:
		logger.GetLogger().Errorf("Unknown storage mode")
		return "", fmt.Errorf("unknown storage mode")
	}
}
