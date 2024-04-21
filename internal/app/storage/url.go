package storage

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type URLData struct {
	ID          string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

var (
	StoragePath string
	UrlsData    []URLData
	mu          sync.Mutex
)

func Init(fileName string) {

	if fileName == "" {
		log.Println("Функция записи на диск отключена")
		return
	}

	StoragePath = fileName
	fmt.Printf("Init FileStoragePath = %v\n", StoragePath)

	err := LoadData(StoragePath)
	if err != nil {
		log.Fatalf("Ошибка при загрузке данных из файла: %v", err)
	}

}

func SaveURLsData() error {
	mu.Lock()
	defer mu.Unlock()

	fmt.Printf("Init FileStoragePath = %v\n", StoragePath)
	if err := ensureDir(StoragePath); err != nil {
		log.Printf("Не удалось создать директорию: %v", err)
		return err
	}

	data, err := json.MarshalIndent(UrlsData, "", " ")
	if err != nil {
		return err
	}

	fmt.Printf("Init WriteFile = %v\n", StoragePath)
	err = os.WriteFile(StoragePath, data, 0644)
	if err != nil {
		log.Printf("Ошибка при записи данных URL в файл: %v", err)
		return err
	}

	return nil
}

func LoadData(dataPath string) error {
	file, err := os.OpenFile(dataPath, os.O_CREATE|os.O_RDWR, 0644)
	fmt.Printf("LoadData file = %v, путь = %v\n", file, dataPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("Файл не найден, инициализация пустого списка URL")
			UrlsData = []URLData{}
			return nil
		}
		log.Printf("Ошибка при открытии файла: %v", err)
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&UrlsData)
	if err != nil {
		log.Printf("Ошибка при декодировании данных из файла: %v", err)
		return err
	}
	return nil
}

func ensureDir(filePath string) error {
	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}
