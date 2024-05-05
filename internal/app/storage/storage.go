package storage

import (
	"fmt"

	"github.com/dualex23/go-url-shortener/internal/app/logger"
	"github.com/google/uuid"
)

func NewStorage(fileName string, mode string, db DataBaseInterface) *Storage {
	storage := &Storage{
		StorageMode: mode,
		StoragePath: fileName,
		DataBase:    db,
	}

	if fileName == "" {
		logger.GetLogger().Error("Disk write functionality is disabled.\n")

		return storage
	}

	if err := storage.LoadData(); err != nil {
		logger.GetLogger().Errorf("Error loading data from file: %v, initializing empty URL list.\n", err)

		storage.UrlsMap = make(map[string]URLData)
		return storage
	}

	return storage
}

// func (s *Storage) Load(fileName string) error {
// 	if fileName == "" {
// 		logger.GetLogger().Error("Disk write functionality is disabled.\n")

// 		return storage
// 	}

// 	if err := storage.LoadData(); err != nil {
// 		logger.GetLogger().Errorf("Error loading data from file: %v, initializing empty URL list.\n", err)

// 		storage.UrlsMap = make(map[string]URLData)
// 		return storage
// 	}
// }

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
		if err := s.DataBase.SaveUrlDB(urlData.ID, urlData.ShortURL, urlData.OriginalURL); err != nil {
			return "", "", err
		}
	case "file":
		s.mu.Lock()
		s.UrlsMap[id] = urlData
		s.mu.Unlock()

		if err := s.SaveURLsFile(); err != nil {
			return "", "", err
		}
	default:
		s.mu.Lock()
		s.UrlsMap[id] = urlData
		s.mu.Unlock()
	}

	return shortURL, id, nil
}

func (s *Storage) FindByID(id string) error {
	switch s.StorageMode {
	case "db":
		// Логика для поиска в базе данных
	case "file":
		// Логика для поиска в файле
	default:
		// Логика для поиска в памяти
	}

	return nil
}

func (s *Storage) Delete(id string) error {
	switch s.StorageMode {
	case "db":
		// Логика для удаления в базе данных
	case "file":
		// Логика для удаления в файле
	default:
		// Логика для удаления в памяти
	}

	return nil
}
