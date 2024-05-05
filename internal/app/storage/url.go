package storage

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/dualex23/go-url-shortener/internal/app/logger"
)

func (s *Storage) SaveURLsFile() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := ensureDir(s.StoragePath); err != nil {
		logger.GetLogger().Errorf("Error creating directory for file: %v\n", err)

		return err
	}

	data, err := json.MarshalIndent(s.UrlsMap, "", " ")
	if err != nil {
		logger.GetLogger().Errorf("Error serializing URL data to JSON: %v\n", err)

		return err
	}

	if err := os.WriteFile(s.StoragePath, data, 0644); err != nil {
		logger.GetLogger().Errorf("Error writing URL data to file: %v\n", err)

		return err
	}

	return nil
}

func (s *Storage) LoadData() error {
	file, err := os.OpenFile(s.StoragePath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		logger.GetLogger().Errorf("Error opening file: %v\n", err)

		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		logger.GetLogger().Errorf("Error getting file information: %v\n", err)

		return err
	}

	if fileInfo.Size() == 0 {
		logger.GetLogger().Error("File is empty, initializing empty URL list.\n")

		s.UrlsMap = make(map[string]URLData)
		return nil
	}

	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&s.UrlsMap); err != nil {
		logger.GetLogger().Errorf("Error decoding data from file: %v\n", err)

		return err
	}

	return nil
}

func ensureDir(filePath string) error {
	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		logger.GetLogger().Infof("Directory does not exist, attempting to create: %s\n", dir)

		if err = os.MkdirAll(dir, 0644); err != nil {
			logger.GetLogger().Errorf("Failed to create directory: %v\n", err)

			return err
		}

		logger.GetLogger().Info("Directory successfully created\n")
	}
	return nil
}
