package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/dualex23/go-url-shortener/internal/app/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DataBase struct {
	DB *sql.DB
}

type DataBaseInterface interface {
	Ping() error
	Close()
	SaveUrls(id, shortURL, originalURL string) error
	LoadUrls() (map[string]URLData, error)
	LoadUrlByID(id string) (*URLData, error)
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

func (db *DataBase) SaveUrls(id, shortURL, originalURL string) error {
	logger.GetLogger().Info("SaveUrls to DB")

	query := `INSERT INTO urls (uuid, short_url, original_url) VALUES ($1,$2,$3)`
	_, err := db.DB.Exec(query, id, shortURL, originalURL)
	if err != nil {
		logger.GetLogger().Errorf("Failed to insert URL: %v", err)
		return err
	}
	return nil
}

func (db *DataBase) LoadUrls() (map[string]URLData, error) {
	logger.GetLogger().Info("LoadUrls from DB")

	urls := make(map[string]URLData)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT uuid, short_url, original_url FROM urls`
	rows, err := db.DB.QueryContext(ctx, query)
	if err != nil {
		logger.GetLogger().Errorf("Failed to execute query in LoadUrls: %v", err)

		return nil, err
	}

	for rows.Next() {
		var u URLData
		if err := rows.Scan(&u.ID, &u.ShortURL, &u.OriginalURL); err != nil {
			logger.GetLogger().Errorf("Failed to scan row: %v", err)
			return nil, err
		}
		urls[u.ID] = u
	}

	if err := rows.Err(); err != nil {
		logger.GetLogger().Errorf("Rows iteration error: %v", err)
		return nil, err
	}

	return urls, nil
}

func (db *DataBase) LoadUrlByID(id string) (*URLData, error) {
	logger.GetLogger().Info("Load url by id")
	var u URLData

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT uuid, short_url, original_url FROM urls WHERE uuid = $1`
	row := db.DB.QueryRowContext(ctx, query, id)

	if err := row.Scan(&u.ID, &u.ShortURL, &u.OriginalURL); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no URL found with ID %s", id)
		}
		logger.GetLogger().Errorf("Failed to execute query in LoadUrlByID: %s", err)
		return nil, err
	}

	return &u, nil
}
