package storage

import "sync"

type Storage struct {
	StorageMode string
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
