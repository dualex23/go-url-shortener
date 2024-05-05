package storage

import (
	"database/sql"

	"github.com/dualex23/go-url-shortener/internal/app/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DataBase struct {
	DB *sql.DB
}

type DataBaseInterface interface {
	Ping() error
	Close()
	SaveUrlDB(id, shortURL, originalURL string) error
}

func NewDB(dataBaseDSN string) (*DataBase, error) {

	db, err := sql.Open("pgx", dataBaseDSN)
	if err != nil {
		logger.GetLogger().Fatalf("Unable to connect to database: %v", err)
	}

	return &DataBase{DB: db}, nil
}

func (db *DataBase) Close() {
	db.DB.Close()
}
func (db *DataBase) Ping() error {
	return db.DB.Ping()
}

func (db *DataBase) SaveUrlDB(id, shortURL, originalURL string) error {
	query := `INSERT INTO urls (uuid, short_url, original_urls) VALUES ($1,$2,$3)`
	_, err := db.DB.Exec(query, id, shortURL, originalURL)
	if err != nil {
		logger.GetLogger().Errorf("Failed to insert URL: %v", err)
		return err
	}
	return nil
}
