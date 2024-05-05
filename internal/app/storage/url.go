package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"github.com/dualex23/go-url-shortener/internal/app/logger"
)

type Storage struct {
	StoragePath string
	UrlsMap     map[string]URLData
	mu          sync.Mutex
	DataBase    DataBaseInterface
}

type URLData struct {
	ID          string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func NewStorage(fileName string, db DataBaseInterface) *Storage {
	storage := &Storage{
		StoragePath: fileName,
		DataBase:    db,
	}
	if fileName == "" {
		logger.GetLogger().Error("Disk write functionality is disabled.")

		return storage
	}

	if err := storage.LoadData(); err != nil {
		logger.GetLogger().Errorf("Error loading data from file: %v, initializing empty URL list.", err)

		storage.UrlsMap = make(map[string]URLData)
		return storage
	}

	return storage
}

func (s *Storage) SaveURLsData() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := ensureDir(s.StoragePath); err != nil {
		logger.GetLogger().Errorf("Error creating directory for file: %v", err)

		return err
	}

	data, err := json.MarshalIndent(s.UrlsMap, "", " ")
	if err != nil {
		logger.GetLogger().Errorf("Error serializing URL data to JSON: %v", err)

		return err
	}

	if err := os.WriteFile(s.StoragePath, data, 0644); err != nil {
		logger.GetLogger().Errorf("Error writing URL data to file: %v", err)

		return err
	}

	return nil
}

func (s *Storage) LoadData() error {
	file, err := os.OpenFile(s.StoragePath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		logger.GetLogger().Errorf("Error opening file: %v", err)

		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		logger.GetLogger().Errorf("Error getting file information: %v", err)

		return err
	}

	if fileInfo.Size() == 0 {
		logger.GetLogger().Error("File is empty, initializing empty URL list.")

		s.UrlsMap = make(map[string]URLData)
		return nil
	}

	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&s.UrlsMap); err != nil {
		logger.GetLogger().Errorf("Error decoding data from file: %v", err)

		return err
	}

	return nil
}

func ensureDir(filePath string) error {
	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		logger.GetLogger().Infof("Directory does not exist, attempting to create: %s", dir)

		if err = os.MkdirAll(dir, 0644); err != nil {
			logger.GetLogger().Errorf("Failed to create directory: %v", err)

			return err
		}

		logger.GetLogger().Info("Directory successfully created")
	}
	return nil
}
