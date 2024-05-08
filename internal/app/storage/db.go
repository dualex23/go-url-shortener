package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
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
	LoadURLByID(id string) (*URLData, error)
	BatchSaveUrls(urls []URLData) error
	FindByOriginalURL(ctx context.Context, originalURL string) (string, string, error)
}

func NewDB(dataBaseDSN string) (*DataBase, error) {
	db, err := sql.Open("pgx", dataBaseDSN)
	if err != nil {
		logger.GetLogger().Fatalf("Unable to connect to database: %v", err)
	}

	dataBase := &DataBase{DB: db}

	if err := dataBase.CreateTable(); err != nil {
		return nil, err
	}

	return dataBase, nil
}

func (db *DataBase) CreateTable() error {
	logger.GetLogger().Info("Starting table and index creation transaction")

	tx, err := db.DB.Begin()
	if err != nil {
		logger.GetLogger().Errorf("Failed to start transaction: %v", err)
		return err
	}

	createTableQuery := `
    CREATE TABLE IF NOT EXISTS urls (
        uuid VARCHAR(255) PRIMARY KEY,
        short_url TEXT NOT NULL,
        original_url TEXT NOT NULL
    );
    `

	if _, err := tx.Exec(createTableQuery); err != nil {
		tx.Rollback()
		logger.GetLogger().Errorf("Failed to create table 'urls': %v", err)
		return err
	}

	createIndexQuery := `CREATE UNIQUE INDEX IF NOT EXISTS idx_original_url ON urls(original_url);`

	if _, err := tx.Exec(createIndexQuery); err != nil {
		tx.Rollback()
		logger.GetLogger().Errorf("Failed to create unique index on original_url: %v", err)
		return err
	}

	if err := tx.Commit(); err != nil {
		logger.GetLogger().Errorf("Failed to commit transaction: %v", err)
		return err
	}

	logger.GetLogger().Info("Table and unique index creation completed successfully")
	return nil
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

func (db *DataBase) LoadURLByID(id string) (*URLData, error) {
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

func (db *DataBase) BatchSaveUrls(urls []URLData) error {

	tx, err := db.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
    INSERT INTO urls (uuid, short_url, original_url)
    VALUES ($1, $2, $3)
    ON CONFLICT (original_url) DO NOTHING
    `)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, url := range urls {
		result, err := stmt.Exec(url.ID, url.ShortURL, url.OriginalURL)
		if err != nil {
			return fmt.Errorf("failed to execute statement: %w", err)
		}
		if affected, _ := result.RowsAffected(); affected == 0 {
			logger.GetLogger().Infoln("No rows inserted due to conflict for URL", url.OriginalURL)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (db *DataBase) FindByOriginalURL(ctx context.Context, originalURL string) (string, string, error) {
	if db.DB == nil {
		logger.GetLogger().Infoln("database connection is not initialized")
		return "", "", fmt.Errorf("database is not initialized")
	}

	var id, shortURL string
	query := `SELECT uuid, short_url FROM urls WHERE original_url = $1`
	err := db.DB.QueryRowContext(ctx, query, originalURL).Scan(&id, &shortURL)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("No entry found for URL: %s", originalURL)
			return "", "", fmt.Errorf("URL not found")
		}
		log.Printf("Error querying database: %v", err)
		return "", "", err
	}

	log.Printf("Found shortened URL for %s: %s", originalURL, shortURL)
	return id, shortURL, nil
}
