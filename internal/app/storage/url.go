package storage

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type Storage struct {
	StoragePath string
	UrlsData    []URLData
	mu          sync.Mutex
}

type URLData struct {
	ID          string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func NewStorage(fileName string) *Storage {
	storage := &Storage{StoragePath: fileName}
	if fileName == "" {
		log.Println("Функция записи на диск отключена.")
		return storage
	}

	if err := storage.LoadData(); err != nil {
		log.Printf("Ошибка при загрузке данных из файла: %v, инициализация пустого списка URL.", err)
		storage.UrlsData = []URLData{}
		return storage
	}

	return storage
}

func (s *Storage) SaveURLsData() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	log.Printf("Попытка записи данных в файл: %s", s.StoragePath)

	if err := ensureDir(s.StoragePath); err != nil {
		log.Printf("Не удалось создать директорию для файла: %v", err)
		return err
	}

	data, err := json.MarshalIndent(s.UrlsData, "", " ")
	if err != nil {
		log.Printf("Ошибка при сериализации данных URL в JSON: %v", err)
		return err
	}

	if err := os.WriteFile(s.StoragePath, data, 0644); err != nil {
		log.Printf("Ошибка при записи данных URL в файл: %v", err)
		return err
	}

	log.Println("Данные успешно сохранены в файл")
	return nil
}

func (s *Storage) LoadData() error {
	file, err := os.OpenFile(s.StoragePath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Printf("Ошибка при открытии файла: %v", err)
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Printf("Ошибка при получении информации о файле: %v", err)
		return err
	}

	if fileInfo.Size() == 0 {
		log.Println("Файл пуст, инициализация пустого списка URL.")
		s.UrlsData = []URLData{}
		return nil
	}

	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&s.UrlsData); err != nil {
		log.Printf("Ошибка при декодировании данных из файла: %v", err)
		return err
	}

	log.Println("Данные успешно загружены из файла.")
	return nil
}

func ensureDir(filePath string) error {
	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Printf("Директория не существует, попытка создать: %s", dir)
		if err = os.MkdirAll(dir, 0755); err != nil {
			log.Printf("Не удалось создать директорию: %v", err)
			return err
		}
		log.Println("Директория успешно создана")
	}
	return nil
}
