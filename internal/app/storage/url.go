package storage

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
)

type URLData struct {
	ID          string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

var (
	StoragePath string = "../../data.json"
	UrlsData    []URLData
	mu          sync.Mutex
)

func Init(fileStoragePath string) {
	fmt.Printf("FileStoragePath = %v\n", fileStoragePath)
	if fileStoragePath == "" {
		log.Println("Функция записи на диск отключена")
		return
	}

	StoragePath = fileStoragePath

	err := LoadData(StoragePath)
	if err != nil {
		log.Fatalf("Ошибка при загрузке данных из файла: %v", err)
	}

}

func SaveURLsData() error {
	mu.Lock()
	defer mu.Unlock()

	data, err := json.MarshalIndent(UrlsData, "", " ")
	if err != nil {
		return err
	}

	err = os.WriteFile(StoragePath, data, 0644)
	if err != nil {
		log.Fatal("Ошибка при записи данных URL в файл:", err)
	}

	return nil
}

func LoadData(dataPath string) error {
	file, err := os.Open(dataPath)
	if err != nil {
		if os.IsNotExist(err) {
			UrlsData = []URLData{}
			return nil
		}
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&UrlsData)
	if err != nil {
		return err
	}
	return nil
}
